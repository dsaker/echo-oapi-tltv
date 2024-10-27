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
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/test"
)

var AudioPauseFilePath = map[int]string{
	3:  "/silence/3SecSilence.mp3",
	4:  "/silence/4SecSilence.mp3",
	5:  "/silence/3SecSilence.mp3",
	6:  "/silence/3SecSilence.mp3",
	7:  "/silence/3SecSilence.mp3",
	8:  "/silence/3SecSilence.mp3",
	9:  "/silence/3SecSilence.mp3",
	10: "/silence/3SecSilence.mp3",
}

func (s *Server) AudioFromFile(e echo.Context) error {
	// get values from multipart form
	titleName := e.FormValue("titleName")
	// convert strings from multipart form to int16's
	fileLangId, err := test.ConvertStringInt16(e.FormValue("fileLanguageId"))
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("error converting fileLanguageId to int16: %s", err.Error()))
	}
	fromLangId, err := test.ConvertStringInt16(e.FormValue("fromLanguageId"))
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("error converting fromLanguageId to int16: %s", err.Error()))
	}
	toLangId, err := test.ConvertStringInt16(e.FormValue("toLanguageId"))
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
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	// get an array of all the phrases from the uploaded file
	stringsSlice, err := audiofile.GetLines(e, src)
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("unable to parse file: %s", err.Error()))
	}
	// TODO add max number of phrases to configs
	if len(stringsSlice) > 100 {
		responseString := fmt.Sprintf("file too large, limit is %d, your file has %d lines", 100, len(stringsSlice))
		return e.String(http.StatusBadRequest, responseString)
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	// insert title into the database
	// TODO roll back on any failure downstream
	title, err := s.queries.InsertTitle(
		e.Request().Context(),
		db.InsertTitleParams{
			Title:        titleName,
			NumSubs:      int16(len(stringsSlice)),
			OgLanguageID: fileLangId,
		})
	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// insert phrases into store as translates object of OgLanguage
	_, err = s.translates.InsertNewPhrases(e, title, s.queries, stringsSlice)
	if err != nil {
		dbErr := s.queries.DeleteTitleById(e.Request().Context(), title.ID)
		if dbErr != nil {
			e.Logger().Error(err)
			return e.String(http.StatusInternalServerError, dbErr.Error())
		}
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	audioFromTitleRequest := oapi.AudioFromTitleJSONRequestBody{
		FromLanguageId: fromLangId,
		TitleId:        title.ID,
		ToLanguageId:   toLangId,
	}
	zipFile, err := createAudioFromTitle(e, s, title, audioFromTitleRequest)
	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.Attachment(zipFile.Name(), title.Title+".zip")
}

func (s *Server) AudioFromTitle(e echo.Context) error {
	var audioFromTitleRequest oapi.AudioFromTitleJSONRequestBody
	err := e.Bind(&audioFromTitleRequest)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	// get title to translate from
	title, err := s.queries.SelectTitleById(e.Request().Context(), audioFromTitleRequest.TitleId)
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

func createAudioFromTitle(e echo.Context, s *Server, t db.Title, r oapi.AudioFromTitleJSONRequestBody) (*os.File, error) {

	// get store.Language for from language from id
	fromLang, err := s.queries.SelectLanguagesById(e.Request().Context(), r.FromLanguageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			e.Logger().Error(err)
			return nil, err
		}
		e.Logger().Error(err)
		return nil, err
	}

	// get store.Language for to language from id
	toLang, err := s.queries.SelectLanguagesById(e.Request().Context(), r.ToLanguageId)
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

	tClient1, err := s.translates.CreateGoogleTranslateClient(e)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}
	ttsClient1, err := s.translates.CreateGoogleTTSClient(e)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}
	// create TTS for from language
	if err = s.translates.CreateTTS(e, s.queries, ttsClient1, tClient1, fromLang, t, fromAudioBasePath); err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	tClient2, err := s.translates.CreateGoogleTranslateClient(e)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}
	ttsClient2, err := s.translates.CreateGoogleTTSClient(e)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}
	// create TTS for to language
	if err = s.translates.CreateTTS(e, s.queries, ttsClient2, tClient2, toLang, t, toAudioBasePath); err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	phraseIds, err := s.queries.SelectPhraseIdsByTitleId(e.Request().Context(), t.ID)
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
	fullPausePath := s.config.TTSBasePath + pausePath

	tmpDirPath := fmt.Sprintf("/tmp/%s-%s/", t.Title, test.RandomString(4))
	err = os.MkdirAll(tmpDirPath, 0777)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	if err = audiofile.BuildAudioInputFiles(e, phraseIds, t, fullPausePath, fromAudioBasePath, toAudioBasePath, tmpDirPath); err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	zipFile, err := audiofile.CreateMp3ZipWithFfmpeg(e, t, tmpDirPath)
	if err != nil {
		return nil, err
	}
	return zipFile, nil
}
