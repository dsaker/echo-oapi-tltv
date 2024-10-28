package audiofile

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	db "talkliketv.click/tltv/db/sqlc"
	audio "talkliketv.click/tltv/internal/audio/pattern"
)

var endSentenceMap = map[rune]bool{
	'!': true,
	'.': true,
	'?': true,
}

type AudioFileX interface {
	GetLines(echo.Context, multipart.File) ([]string, error)
	CreateMp3ZipWithFfmpeg(echo.Context, db.Title, string) (*os.File, error)
	BuildAudioInputFiles(echo.Context, []int64, db.Title, string, string, string, string) error
}

type AudioFile struct {
}

func (a *AudioFile) GetLines(e echo.Context, f multipart.File) ([]string, error) {

	// get file type, options are srt, single line text or paragraph
	fileType := ""
	scanner := bufio.NewScanner(f)
	count := 0
	var line string

	// verify if file is srt
	for scanner.Scan() {
		if fileType != "" || count > 4 {
			break
		}
		line = scanner.Text()
		// if line contains ">" and doesn't contain any letters it is srt file
		if strings.Contains(line, ">") {
			containsAlpha, err := regexp.MatchString("[a-zA-Z]", line)
			if err != nil {
				e.Logger().Error(err)
				return nil, err
			}
			if !containsAlpha {
				fileType = "srt"
			}
		}
		count++
	}
	//start at the first line again
	_, err := f.Seek(0, 0)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}
	count = 0
	scanner = bufio.NewScanner(f)
	// verify if file is in paragraph form
	for scanner.Scan() {
		if fileType != "" || count > 4 {
			break
		}
		line = scanner.Text()
		// Split on punctuation characters
		re := regexp.MustCompile(`[.!?]`)
		result := re.Split(line, -1)
		if len(result) > 3 {
			fileType = "paragraph"
		}
		count++
	}
	// TODO somehow verify single phrase per line form (these can be multiple sentences
	_, err = f.Seek(0, 0)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}
	var stringsSlice []string
	if fileType == "srt" {
		stringsSlice = parseSrt(f)
	}
	if fileType == "paragraph" {
		stringsSlice = parseParagraph(f)
	}
	if fileType == "" {
		stringsSlice = parseSingle(f)
	}
	if len(stringsSlice) == 0 {
		return nil, errors.New("unable to parse file")
	}

	return stringsSlice, nil
}

func parseSrt(f multipart.File) []string {
	var stringsSlice []string
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	var line string
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		} else if line[0] >= '0' && line[0] <= '9' {
			continue
		} else if line[0] == '[' && line[len(line)-1] == ']' {
			continue
		} else {
			// if the next line following subtitle is not new line it is more dialogue so combine it
			scanner.Scan()
			nextLine := scanner.Text()
			if nextLine != "" {
				line = strings.ReplaceAll(line, "\n", "")
				line = line + " " + nextLine
				line = replaceFmt(line)
			} else {
				line = replaceFmt(line)
			}
		}

		// if sentence is too short don't keep it
		words := strings.Fields(line)
		if len(words) > 3 {
			stringsSlice = append(stringsSlice, line)
		}
	}

	return stringsSlice
}

func parseParagraph(f multipart.File) []string {
	var stringsSlice []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Split on punctuation characters
		last := 0
		for i, c := range line {
			if i == len(line)-1 {
				sentence := strings.TrimSpace(line[last+1 : i+1])
				fmt.Println(sentence)
				last = i
				words := strings.Fields(sentence)
				if len(words) > 3 {
					stringsSlice = append(stringsSlice, line)
				}
			} else if endSentenceMap[c] {
				sentence := strings.TrimSpace(line[last+1 : i+1])
				fmt.Println(sentence)
				last = i
				words := strings.Fields(sentence)
				if len(words) > 3 {
					stringsSlice = append(stringsSlice, line)
				}
			}
		}
	}

	return stringsSlice
}

func parseSingle(f multipart.File) []string {
	var stringsSlice []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		if len(words) > 2 {
			stringsSlice = append(stringsSlice, line)
		}
	}

	return stringsSlice
}

func replaceFmt(line string) string {
	// remove any characters between brackets and brackets [...] or {...} or <...>
	re := regexp.MustCompile("\\[.*?]")
	line = re.ReplaceAllString(line, "")
	re = regexp.MustCompile("\\{.*?}")
	line = re.ReplaceAllString(line, "")
	re = regexp.MustCompile("<.*?>")
	line = strings.ReplaceAll(line, "-", "")
	line = strings.ReplaceAll(line, "\"", "")
	line = strings.ReplaceAll(line, "'", "")

	return line
}

func (a *AudioFile) CreateMp3ZipWithFfmpeg(e echo.Context, t db.Title, tmpDir string) (*os.File, error) {
	// get a list of files from the temp directory
	files, err := os.ReadDir(tmpDir)
	// create outputs folder to hold all the mp3's to zip
	outDirPath := tmpDir + "outputs"
	err = os.MkdirAll(outDirPath, 0777)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}
	for i, f := range files {
		//ffmpeg -f concat -safe 0 -i ffmpeg_input.txt -c copy output.mp3
		outputString := fmt.Sprintf("%s/%s-%d.mp3", outDirPath, t.Title, i)
		// TODO make this concurrent
		cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", tmpDir+f.Name(), "-c", "copy", outputString)

		//Execute the command and get the output
		output, err := cmd.CombinedOutput()
		if err != nil {
			e.Logger().Error(err)
			e.Logger().Error(string(output))
			return nil, err
		}
	}

	zipFile, err := os.Create(tmpDir + "/" + t.Title + ".zip")
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// get a list of files from the output directory
	files, err = os.ReadDir(outDirPath)
	for _, file := range files {
		err = addFileToZip(e, zipWriter, outDirPath+"/"+file.Name())
		if err != nil {
			return nil, err
		}
	}

	return zipFile, err
}

func addFileToZip(e echo.Context, zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		e.Logger().Error(err)
		return err
	}
	defer file.Close()

	fInfo, err := file.Stat()
	if err != nil {
		e.Logger().Error(err)
		return err
	}

	header, err := zip.FileInfoHeader(fInfo)
	if err != nil {
		e.Logger().Error(err)
		return err
	}

	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		e.Logger().Error(err)
		return err
	}

	_, err = io.Copy(writer, file)
	e.Logger().Info("wrote file: %s", file.Name())
	return err
}

// TODO shorten the parameters
func (a *AudioFile) BuildAudioInputFiles(e echo.Context, ids []int64, t db.Title, pause, from, to, tmpDir string) error {

	pMap := make(map[int]int64)

	// map phrase ids to zero through len(phrase ids) to map correctly to pattern.Pattern
	for i, pid := range ids {
		pMap[i] = pid
	}

	maxP := slices.Max(ids)
	// create chunks of []Audio pattern to split up audio files into ~20 minute lengths
	// TODO look at slices.Chunk to see how it accepts any type of slice
	chunkedSlice := slices.Chunk(audio.Pattern, 250)
	count := 1
	last := false
	for chunk := range chunkedSlice {
		inputString := fmt.Sprintf("%s-input-%d", t.Title, count)
		count++
		f, err := os.Create(tmpDir + inputString)
		if err != nil {
			e.Logger().Error(err)
			return err
		}
		defer f.Close()

		for _, audioStruct := range chunk {
			// if: we have reached the highest phrase id then this will be the last audio block
			// this will also skip non-existent phrase ids
			// else if: native language then we add filepath for from audio mp3
			// else: add audio filepath for language you want to learn
			phraseId := pMap[audioStruct.Id]
			if phraseId == maxP {
				last = true
			} else if phraseId == 0 && audioStruct.Id > 0 {
				continue
			} else if audioStruct.Native == true {
				_, err = f.WriteString(fmt.Sprintf("file '%s%d'\n", from, phraseId))
				_, err = f.WriteString(fmt.Sprintf("file '%s'\n", pause))
				if err != nil {
					e.Logger().Error(err)
					return err
				}
			} else {
				_, err = f.WriteString(fmt.Sprintf("file '%s%d'\n", to, phraseId))
				_, err = f.WriteString(fmt.Sprintf("file '%s'\n", pause))
				if err != nil {
					e.Logger().Error(err)
					return err
				}
			}
		}
		if last {
			break
		}
	}

	return nil
}
