package api

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/hyacinthus/mp3join"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"strconv"
	db "talkliketv.click/tltv/db/sqlc"
)

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
		strconv.Itoa(int(title.ID)) + "/" +
		strconv.Itoa(int(audioFromTitleRequest.FromLanguageId)) + "/"

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

	destPath := audioBasePath + "/" + "output.mp3"
	if err = buildAudioFile(phraseIds, "3", destPath, fromAudioBasePath, toAudioBasePath); err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}
	return eCtx.String(http.StatusOK, "OK")
}

func buildAudioFile(ids []int64, destPath, pauseL, from, to string) error {

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	pause, err := os.ReadFile(dir + "/internal/silence/" + pauseL)
	if err != nil {
		return err
	}
	pauseReader := bytes.NewReader(pause)

	// https://github.com/hyacinthus/mp3join/blob/master/joiner.go
	joiner := mp3join.New()

	// readers is the input mp3 files
	for i, pid := range ids {
		stringPid := strconv.FormatInt(pid, 10)
		// get from mp3 from directory
		fiFrom, err := os.ReadFile(from + stringPid)
		if err != nil {
			return err
		}

		// add one from phrase with pause
		err = addMp3withPause(fiFrom, pauseReader, joiner)
		if err != nil {
			return err
		}

		fiTo, err := os.ReadFile(to + stringPid)
		if err != nil {
			return err
		}
		// add two to phrases with pause
		err = addMp3withPause(fiTo, pauseReader, joiner)
		if err != nil {
			return err
		}
		err = addMp3withPause(fiTo, pauseReader, joiner)
		if err != nil {
			return err
		}
		if i == 30 {
			break
		}
	}

	dest := joiner.Reader()
	// Create the output file
	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the contents of the reader to the file
	_, err = io.Copy(file, dest)
	if err != nil {
		return err
	}
	return nil
}

func addMp3withPause(b []byte, p *bytes.Reader, j *mp3join.Joiner) error {
	// convert to reader and append
	reader := bytes.NewReader(b)
	err := j.Append(reader)
	if err != nil {
		return err
	}

	// add pause
	err = j.Append(p)
	if err != nil {
		return err
	}

	return nil
}

func createTTS(eCtx echo.Context, s *Server, lang db.Language, title db.Title, basepath string) error {
	skip, err := pathExists(basepath)
	if err != nil {
		return err
	}

	if !skip {
		fromTranslates, err := getOrCreateTranslates(eCtx, s, title.ID, lang, title.OgLanguageID)
		if err != nil {
			return err
		}

		if err = s.translates.TextToSpeech(eCtx, fromTranslates, basepath, lang.Tag); err != nil {
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
		return nil, err
	}

	// create translates for title and to language and return
	translatesReturn, err := s.translates.TranslatePhrases(eCtx, fromTranslates, toLang)
	if err != nil {
		return nil, err
	}

	dbTranslates, err := s.translates.InsertTranslates(eCtx, s.queries, toLang.ID, translatesReturn)
	if err != nil {
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
