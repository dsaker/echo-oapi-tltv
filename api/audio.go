package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/audio/audiofile"
	"talkliketv.click/tltv/internal/util"
)

var AudioPauseFilePath = map[int]string{

	3:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/audio/silence/3SecSilence.mp3",
	4:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/audio/silence/4SecSilence.mp3",
	5:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/audio/silence/3SecSilence.mp3",
	6:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/audio/silence/3SecSilence.mp3",
	7:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/audio/silence/3SecSilence.mp3",
	8:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/audio/silence/3SecSilence.mp3",
	9:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/audio/silence/3SecSilence.mp3",
	10: "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/audio/silence/3SecSilence.mp3",
}

func (s *Server) AudioFromFile(e echo.Context) error {
	// get values from multipart form
	titleName := e.FormValue("titleName")
	// convert strings from multipart form to int16's
	fileLangId, err := util.ConvertStringInt16(e.FormValue("fileLanguageId"))
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("error converting fileLanguageId to int16: %s", err.Error()))
	}
	fromLangId, err := util.ConvertStringInt16(e.FormValue("fromLanguageId"))
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("error converting fromLanguageId to int16: %s", err.Error()))
	}
	toLangId, err := util.ConvertStringInt16(e.FormValue("toLanguageId"))
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("error converting toLanguageId to int16: %s", err.Error()))
	}

	// Get file handler for filename, size and headers
	file, err := e.FormFile("filePath")
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	// Check if file size is too large 64000 == 8KB ~ approximately 4 pages of text
	if file.Size > s.config.FileUploadLimit*8000 {
		rString := fmt.Sprintf("file too large (%d > %d)", file.Size, s.config.FileUploadLimit*8000)
		return e.String(http.StatusBadRequest, rString)
	}
	src, err := file.Open()
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	// get an array of all the phrases from the uploaded file
	stringsSlice, err := audiofile.GetLines(e, src)
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("unable to parse file: %s", err.Error()))
	}
	// TODO add max number of phrases to configs
	if len(stringsSlice) > 800 {
		responseString := fmt.Sprintf("file too large, limit is %d, your file has %d lines", 800, len(stringsSlice))
		return e.String(http.StatusBadRequest, responseString)
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	// insert title into the database
	// TODO roll back on any failure downstream
	title, err := s.Queries.InsertTitle(
		e.Request().Context(),
		db.InsertTitleParams{
			Title:        titleName,
			NumSubs:      int16(len(stringsSlice)),
			OgLanguageID: fileLangId,
		})
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// insert phrases into db as translates object of OgLanguage
	_, err = s.Translates.InsertNewPhrases(e, title, s.Queries, stringsSlice)
	if err != nil {
		dbErr := s.Queries.DeleteTitleById(e.Request().Context(), title.ID)
		if dbErr != nil {
			return e.String(http.StatusInternalServerError, dbErr.Error())
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	audioFromTitleRequest := AudioFromTitleJSONRequestBody{
		FromLanguageId: fromLangId,
		TitleId:        title.ID,
		ToLanguageId:   toLangId,
	}
	zipFile, err := createAudioFromTitle(e, s, title, audioFromTitleRequest)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.Attachment(zipFile.Name(), title.Title+".zip")
}

func (s *Server) AudioFromTitle(e echo.Context) error {
	var audioFromTitleRequest AudioFromTitleJSONRequestBody
	err := e.Bind(&audioFromTitleRequest)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	// get title to translate from
	title, err := s.Queries.SelectTitleById(e.Request().Context(), audioFromTitleRequest.TitleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.String(http.StatusBadRequest, "invalid title id")
		}
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	zipFile, err := createAudioFromTitle(e, s, title, audioFromTitleRequest)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.Attachment(zipFile.Name(), title.Title+".zip")
}

func createAudioFromTitle(e echo.Context, s *Server, t db.Title, r AudioFromTitleJSONRequestBody) (*os.File, error) {

	// get db.Language for from language from id
	fromLang, err := s.Queries.SelectLanguagesById(e.Request().Context(), r.FromLanguageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			e.Logger().Error(err)
			return nil, err
		}
		e.Logger().Error(err)
		return nil, err
	}

	// get db.Language for to language from id
	toLang, err := s.Queries.SelectLanguagesById(e.Request().Context(), r.ToLanguageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			e.Logger().Error(err)
			return nil, err
		}
		e.Logger().Error(err)
		return nil, err
	}

	audioBasePath := fmt.Sprintf("%s%d/", s.config.TTSBasePath, t.ID)

	fromAudioBasePath := fmt.Sprintf("%s%d/", audioBasePath, r.FromLanguageId)
	toAudioBasePath := fmt.Sprintf("%s%d/", audioBasePath, r.ToLanguageId)

	if err = audiofile.CreateTTS(e, s, fromLang, t, fromAudioBasePath); err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	if err = audiofile.CreateTTS(e, s, toLang, t, toAudioBasePath); err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	phraseIds, err := s.Queries.SelectPhraseIdsByTitleId(e.Request().Context(), t.ID)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	// get pause audio file for length provided by user
	// TODO add pause length to configs
	pausePath, ok := AudioPauseFilePath[4]
	if !ok {
		e.Logger().Error(err)
		return nil, err
	}

	tmpDirPath := fmt.Sprintf("/tmp/%s-%s/", t.Title, util.RandomString(4))
	err = os.MkdirAll(tmpDirPath, 0777)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}
	defer os.RemoveAll(tmpDirPath)

	if err = audiofile.BuildAudioInputFiles(e, phraseIds, t, pausePath, fromAudioBasePath, toAudioBasePath, tmpDirPath); err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	zipFile, err := audiofile.CreateMp3ZipWithFfmpeg(e, t, tmpDirPath)
	if err != nil {
		return nil, err
	}
	return zipFile, nil
}

//func createMp3ZipWithFfmpeg(e echo.Context, t db.Title, tmpDir string) (*os.File, error) {
//	// get a list of files from the temp directory
//	files, err := os.ReadDir(tmpDir)
//	// create outputs folder to hold all the mp3's to zip
//	outDirPath := tmpDir + "outputs"
//	err = os.MkdirAll(outDirPath, 0777)
//	if err != nil {
//		e.Logger().Error(err)
//		return nil, err
//	}
//	for i, f := range files {
//		//ffmpeg -f concat -safe 0 -i ffmpeg_input.txt -c copy output.mp3
//		outputString := fmt.Sprintf("%s/%s-%d.mp3", outDirPath, t.Title, i)
//		cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", tmpDir+f.Name(), "-c", "copy", outputString)
//
//		//Execute the command and get the output
//		output, err := cmd.CombinedOutput()
//		if err != nil {
//			e.Logger().Error(err)
//			e.Logger().Error(string(output))
//			return nil, err
//		}
//	}
//
//	zipFile, err := os.Create(tmpDir + "/" + t.Title + ".zip")
//	if err != nil {
//		e.Logger().Error(err)
//		return nil, err
//	}
//	defer zipFile.Close()
//
//	zipWriter := zip.NewWriter(zipFile)
//	defer zipWriter.Close()
//
//	// get a list of files from the output directory
//	files, err = os.ReadDir(outDirPath)
//	for _, file := range files {
//		err = addFileToZip(e, zipWriter, outDirPath+"/"+file.Name())
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return zipFile, err
//}
//
//func addFileToZip(e echo.Context, zipWriter *zip.Writer, filename string) error {
//	file, err := os.Open(filename)
//	if err != nil {
//		e.Logger().Error(err)
//		return err
//	}
//	defer file.Close()
//
//	fInfo, err := file.Stat()
//	if err != nil {
//		e.Logger().Error(err)
//		return err
//	}
//
//	header, err := zip.FileInfoHeader(fInfo)
//	if err != nil {
//		e.Logger().Error(err)
//		return err
//	}
//
//	header.Name = filepath.Base(filename)
//	header.Method = zip.Deflate
//
//	writer, err := zipWriter.CreateHeader(header)
//	if err != nil {
//		e.Logger().Error(err)
//		return err
//	}
//
//	_, err = io.Copy(writer, file)
//	e.Logger().Info("wrote file: %s", file.Name())
//	return err
//}
//
//func buildAudioInputFiles(e echo.Context, ids []int64, t db.Title, pause, from, to, tmpDir string) error {
//
//	pMap := make(map[int]int64)
//
//	// map phrase ids to zero through len(phrase ids) to map correctly to pattern.Pattern
//	for i, pid := range ids {
//		pMap[i] = pid
//	}
//
//	maxP := slices.Max(ids)
//	// create chunks of []Audio pattern to split up audio files into ~20 minute lengths
//	// TODO look at slices.Chunk to see how it accepts any type of slice
//	chunkedSlice := slices.Chunk(audio.Pattern, 250)
//	count := 1
//	last := false
//	for chunk := range chunkedSlice {
//		inputString := fmt.Sprintf("%s-input-%d", t.Title, count)
//		count++
//		f, err := os.Create(tmpDir + inputString)
//		if err != nil {
//			e.Logger().Error(err)
//			return err
//		}
//		defer f.Close()
//
//		for _, audioStruct := range chunk {
//			// if: we have reached the highest phrase id then this will be the last audio block
//			// this will also skip non-existent phrase ids
//			// else if: native language then we add filepath for from audio mp3
//			// else: add audio filepath for language you want to learn
//			phraseId := pMap[audioStruct.Id]
//			if phraseId == maxP {
//				last = true
//			} else if phraseId == 0 && audioStruct.Id > 0 {
//				continue
//			} else if audioStruct.Native == true {
//				_, err = f.WriteString(fmt.Sprintf("file '%s%d'\n", from, phraseId))
//				_, err = f.WriteString(fmt.Sprintf("file '%s'\n", pause))
//				if err != nil {
//					e.Logger().Error(err)
//					return err
//				}
//			} else {
//				_, err = f.WriteString(fmt.Sprintf("file '%s%d'\n", to, phraseId))
//				_, err = f.WriteString(fmt.Sprintf("file '%s'\n", pause))
//				if err != nil {
//					e.Logger().Error(err)
//					return err
//				}
//			}
//		}
//		if last {
//			break
//		}
//	}
//
//	return nil
//}
//
//func createTTS(eCtx echo.Context, s *Server, lang db.Language, title db.Title, basepath string) error {
//	skip, err := pathExists(basepath)
//	if err != nil {
//		eCtx.Logger().Error(err)
//		return err
//	}
//
//	if !skip {
//		fromTranslates, err := getOrCreateTranslates(eCtx, s, title.ID, lang, title.OgLanguageID)
//		if err != nil {
//			return err
//		}
//
//		err = os.MkdirAll(basepath, 0777)
//		if err != nil {
//			eCtx.Logger().Error(err)
//			return err
//		}
//
//		if err = s.translates.TextToSpeech(eCtx, fromTranslates, basepath, lang.Tag); err != nil {
//			eCtx.Logger().Error(err)
//			return err
//		}
//	}
//
//	return nil
//}
//
//func getOrCreateTranslates(eCtx echo.Context, s *Server, titleId int64, toLang db.Language, fromLangId int16) ([]db.Translate, error) {
//	// see if translates exist for title for language
//	exists, err := s.queries.SelectExistsTranslates(
//		eCtx.Request().Context(),
//		db.SelectExistsTranslatesParams{
//			LanguageID: toLang.ID,
//			ID:         titleId,
//		})
//
//	// if exists get translates for language
//	if exists {
//		params := db.SelectTranslatesByTitleIdLangIdParams{
//			LanguageID: toLang.ID,
//			ID:         titleId,
//		}
//		translates, err := s.queries.SelectTranslatesByTitleIdLangId(eCtx.Request().Context(), params)
//		if err != nil {
//			eCtx.Logger().Error(err)
//			return nil, err
//		}
//		return translates, nil
//	}
//
//	// if not exists get translates for fromLangId
//	fromTranslates, err := s.queries.SelectTranslatesByTitleIdLangId(
//		eCtx.Request().Context(),
//		db.SelectTranslatesByTitleIdLangIdParams{
//			LanguageID: fromLangId,
//			ID:         titleId,
//		})
//	if err != nil {
//		eCtx.Logger().Error(err)
//		return nil, err
//	}
//
//	// create translates for title and to language and return
//	translatesReturn, err := s.translates.TranslatePhrases(eCtx, fromTranslates, toLang)
//	if err != nil {
//		eCtx.Logger().Error(err)
//		return nil, err
//	}
//
//	dbTranslates, err := s.translates.InsertTranslates(eCtx, s.queries, toLang.ID, translatesReturn)
//	if err != nil {
//		eCtx.Logger().Error(err)
//		return nil, err
//	}
//	return dbTranslates, nil
//}
//
//// pathExists returns whether the given file or directory exists
//func pathExists(path string) (bool, error) {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true, nil
//	}
//	if os.IsNotExist(err) {
//		return false, nil
//	}
//	return false, err
//}
