package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"slices"
	"strconv"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/audio/pattern"
)

// Enum value maps for audio pause filepaths.
type Pause string

var AudioPauseFilePath = map[int32]string{
	3:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/internal/silence/3SecSilence.mp3",
	4:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/internal/silence/4SecSilence.mp3",
	5:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/internal/silence/3SecSilence.mp3",
	6:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/internal/silence/3SecSilence.mp3",
	7:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/internal/silence/3SecSilence.mp3",
	8:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/internal/silence/3SecSilence.mp3",
	9:  "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/internal/silence/3SecSilence.mp3",
	10: "/Users/dustysaker/go/src/github.com/dsaker/echo-oapi-tltv/internal/silence/3SecSilence.mp3",
}

func (s *Server) AudioFromFile(ctx echo.Context) error {
	ctx.Logger().Info("inside AudioFromFile")
	return ctx.NoContent(http.StatusNotImplemented)
}

func (s *Server) AudioFromTitle(eCtx echo.Context) error {
	var audioFromTitleRequest AudioFromTitleJSONRequestBody
	err := eCtx.Bind(&audioFromTitleRequest)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	// get title to translate from
	title, err := s.queries.SelectTitleById(eCtx.Request().Context(), audioFromTitleRequest.TitleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return eCtx.String(http.StatusBadRequest, "invalid title id")
		}
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	// get db.Language for from language from id
	fromLang, err := s.queries.SelectLanguagesById(eCtx.Request().Context(), audioFromTitleRequest.FromLanguageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return eCtx.String(http.StatusBadRequest, "invalid from language id")
		}
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	// get db.Language for to language from id
	toLang, err := s.queries.SelectLanguagesById(eCtx.Request().Context(), audioFromTitleRequest.ToLanguageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return eCtx.String(http.StatusBadRequest, "invalid from language id")
		}
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	audioBasePath := s.config.TTSBasePath +
		strconv.Itoa(int(title.ID)) + "/"

	fromAudioBasePath := audioBasePath + strconv.Itoa(int(audioFromTitleRequest.FromLanguageId)) + "/"
	toAudioBasePath := audioBasePath + strconv.Itoa(int(audioFromTitleRequest.ToLanguageId)) + "/"

	if err = createTTS(eCtx, s, fromLang, title, fromAudioBasePath); err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	if err = createTTS(eCtx, s, toLang, title, toAudioBasePath); err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	phraseIds, err := s.queries.SelectPhraseIdsByTitleId(eCtx.Request().Context(), title.ID)
	if err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	// get pause audio file for length provided by user
	// TODO add pause length to configs
	pausePath, ok := AudioPauseFilePath[4]
	if !ok {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	if err = buildAudioFile(eCtx, phraseIds, title, pausePath, fromAudioBasePath, toAudioBasePath); err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}
	return eCtx.String(http.StatusOK, "OK")
}

func buildAudioFile(e echo.Context, ids []int64, t db.Title, pause, from, to string) error {

	//f, err := os.Create("/tmp/" + "input-" + util.RandomString(6))
	//if err != nil {
	//	return err
	//}
	//defer f.Close()

	pMap := make(map[int]int64)

	// map phrase ids to zero through len(phrase ids) to map correctly to pattern.Pattern
	for i, pid := range ids {
		pMap[i] = pid
	}

	maxP := slices.Max(ids)
	tmpDirString := "/tmp/" + t.Title
	err := os.MkdirAll(tmpDirString, os.ModePerm)
	if err != nil {
		e.Logger().Error(err)
		return err
	}
	// create chunks of []Audio pattern to split up audio files into ~20 minute lengths
	chunkedSlice := slices.Chunk(pattern.Pattern, 250)
	count := 1
	last := false
	for chunk := range chunkedSlice {
		inputString := fmt.Sprintf("%s-input-%d", t.Title, count)
		count++
		f, err := os.Create("/tmp/" + inputString)
		if err != nil {
			e.Logger().Error(err)
			return err
		}
		defer f.Close()

		for _, audio := range chunk {
			// if: we have reached the highest phrase id then this will be the last audio block
			// this will also skip non-existent phrase ids
			// else if: native language then we add filepath for from audio mp3
			// else: add audio filepath for language you want to learn
			phraseId := pMap[audio.Id]
			if phraseId == maxP {
				last = true
			} else if phraseId == 0 && audio.Id > 0 {
				continue
			} else if audio.Native == true {
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
	//for _, pid := range ids {
	//	stringPid := strconv.FormatInt(pid, 10)
	//	_, err = f.WriteString(fmt.Sprintf("file '%s'\n", from+stringPid))
	//	_, err = f.WriteString(fmt.Sprintf("file '%s'\n", pause))
	//	_, err = f.WriteString(fmt.Sprintf("file '%s'\n", to+stringPid))
	//	_, err = f.WriteString(fmt.Sprintf("file '%s'\n", pause))
	//	_, err = f.WriteString(fmt.Sprintf("file '%s'\n", to+stringPid))
	//	_, err = f.WriteString(fmt.Sprintf("file '%s'\n", pause))
	//}

	// ffmpeg -f concat -safe 0 -i ffmpeg_input.txt -c copy output.mp3
	//cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", f.Name(), "-c", "copy", "/tmp/"+t.Title+".mp3")

	// Execute the command and get the output
	//output, err := cmd.CombinedOutput()
	//if err != nil {
	//	return err
	//}
	//
	//if strings.Contains(string(output), "Error") {
	//	return errors.New(string(output))
	//}

	e.Logger().Info(fmt.Sprintf("Created mp3 audio file: %s", "tmp/"+t.Title+".mp3"))

	return nil
}

func createTTS(eCtx echo.Context, s *Server, lang db.Language, title db.Title, basepath string) error {
	skip, err := pathExists(basepath)
	if err != nil {
		eCtx.Logger().Error(err)
		return err
	}

	if !skip {
		fromTranslates, err := getOrCreateTranslates(eCtx, s, title.ID, lang, title.OgLanguageID)
		if err != nil {
			eCtx.Logger().Error(err)
			return err
		}

		err = os.MkdirAll(basepath, 0777)
		if err != nil {
			eCtx.Logger().Error(err)
			return err
		}

		if err = s.translates.TextToSpeech(eCtx, fromTranslates, basepath, lang.Tag); err != nil {
			eCtx.Logger().Error(err)
			return err
		}
	}

	return nil
}

func getOrCreateTranslates(eCtx echo.Context, s *Server, titleId int64, toLang db.Language, fromLangId int16) ([]db.Translate, error) {
	// see if translates exist for title for language
	exists, err := s.queries.SelectExistsTranslates(
		eCtx.Request().Context(),
		db.SelectExistsTranslatesParams{
			LanguageID: toLang.ID,
			ID:         titleId,
		})

	// if exists get translates for language
	if exists {
		params := db.SelectTranslatesByTitleIdLangIdParams{
			LanguageID: toLang.ID,
			ID:         titleId,
		}
		translates, err := s.queries.SelectTranslatesByTitleIdLangId(eCtx.Request().Context(), params)
		if err != nil {
			eCtx.Logger().Error(err)
			return nil, err
		}
		return translates, nil
	}

	// if not exists get translates for fromLangId
	fromTranslates, err := s.queries.SelectTranslatesByTitleIdLangId(
		eCtx.Request().Context(),
		db.SelectTranslatesByTitleIdLangIdParams{
			LanguageID: fromLangId,
			ID:         titleId,
		})
	if err != nil {
		eCtx.Logger().Error(err)
		return nil, err
	}

	// create translates for title and to language and return
	translatesReturn, err := s.translates.TranslatePhrases(eCtx, fromTranslates, toLang)
	if err != nil {
		eCtx.Logger().Error(err)
		return nil, err
	}

	dbTranslates, err := s.translates.InsertTranslates(eCtx, s.queries, toLang.ID, translatesReturn)
	if err != nil {
		eCtx.Logger().Error(err)
		return nil, err
	}
	return dbTranslates, nil
}

// pathExists returns whether the given file or directory exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
