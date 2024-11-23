package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/audio/audiofile"
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/test"
	"talkliketv.click/tltv/internal/util"
)

// AudioFromFile accepts a file in srt, phrase per line, or paragraph form and
// sends a zip file of mp3 audio tracks for learning a language that you choose
func (s *Server) AudioFromFile(e echo.Context) error {
	// get values from multipart form
	titleName := e.FormValue("titleName")
	// convert strings from multipart form to int16's
	fileLangId, err := util.ConvertStringInt16(e.FormValue("fileLanguageId"))
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("error converting fileLanguageId to int16: %s", err.Error()))
	}
	toVoiceId, err := util.ConvertStringInt16(e.FormValue("toVoiceId"))
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("error converting toVoiceId to int16: %s", err.Error()))
	}
	fromVoiceId, err := util.ConvertStringInt16(e.FormValue("fromVoiceId"))
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("error converting fromVoiceId to int16: %s", err.Error()))
	}
	// check if user sent 'pause' in the request and update config if they did
	pause := e.FormValue("pause")
	if pause != "" {
		pauseInt, err := strconv.Atoi(pause)
		if err != nil {
			return e.String(http.StatusBadRequest, fmt.Sprintf("error converting fromVoiceId to int: %s", err.Error()))
		}
		if pauseInt > 10 || pauseInt < 3 {
			return e.String(http.StatusBadRequest, fmt.Sprintf("pause must be between 3 and 10: %d", pauseInt))
		}
		s.config.PhrasePause = pauseInt
	}

	pattern := e.FormValue("pattern")
	if pattern != "" {
		patternInt, err := strconv.Atoi(pattern)
		if err != nil {
			return e.String(http.StatusBadRequest, fmt.Sprintf("error converting pattern to int: %s", err.Error()))
		}
		if patternInt > 3 || patternInt < 1 {
			return e.String(http.StatusBadRequest, fmt.Sprintf("pattern must be between 1 and 3: %d", patternInt))
		}
		s.config.AudioPattern = patternInt
	}
	e.Set("pattern", s.config.AudioPattern)

	title, phraseZipFile, err := s.processFile(e, titleName, fileLangId)
	if err != nil {
		if errors.Is(err, util.ErrTooManyPhrases) {
			return e.Attachment(phraseZipFile.Name(), "TooManyPhrasesUseTheseFiles.zip")
		}
		if strings.Contains(err.Error(), "unable to parse file") {
			return e.String(http.StatusBadRequest, err.Error())
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	audioFromTitleRequest := oapi.AudioFromTitleJSONRequestBody{
		TitleId:     title.ID,
		ToVoiceId:   toVoiceId,
		FromVoiceId: fromVoiceId,
	}
	zipFile, err := s.createAudioFromTitle(e, *title, audioFromTitleRequest)
	if err != nil {
		if errors.Is(err, util.ErrVoiceLangIdNoMatch) {
			return e.String(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, util.ErrVoiceIdInvalid) {
			return e.String(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, util.ErrOneFile) {
			return e.Attachment(zipFile.Name(), title.Title+".mp3")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.Attachment(zipFile.Name(), title.Title+".zip")
}

func (s *Server) processFile(e echo.Context, titleName string, fileLangId int16) (*db.Title, *os.File, error) {
	// Get file handler for filename, size and headers
	fh, err := e.FormFile("filePath")
	if err != nil {
		e.Logger().Error(err)
		return nil, nil, util.ErrUnableToParseFile(err)
	}

	// Check if file size is too large 64000 == 8KB ~ approximately 4 pages of text
	if fh.Size > s.config.FileUploadLimit {
		rString := fmt.Sprintf("file too large (%d > %d)", fh.Size, s.config.FileUploadLimit)
		return nil, nil, util.ErrUnableToParseFile(errors.New(rString))
	}
	src, err := fh.Open()
	if err != nil {
		e.Logger().Error(err)
		return nil, nil, err
	}
	defer src.Close()

	// get an array of all the phrases from the uploaded file
	stringsSlice, err := s.af.GetLines(e, src)
	if err != nil {
		return nil, nil, util.ErrUnableToParseFile(err)
	}
	// send back zip of split files of phrase that requester can use if too big
	if len(stringsSlice) > s.config.MaxNumPhrases {
		chunkedPhrases := slices.Chunk(stringsSlice, s.config.MaxNumPhrases)
		phrasesBasePath := s.config.TTSBasePath + titleName + "/"
		// create zip of phrases files of maxNumPhrases for user to use instead of uploaded file
		zipFile, err := s.af.CreatePhrasesZip(e, chunkedPhrases, phrasesBasePath, titleName)
		if err != nil {
			e.Logger().Error(err)
			return nil, nil, err
		}
		return nil, zipFile, util.ErrTooManyPhrases
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	// insert title into the database
	title, err := s.queries.InsertTitle(
		e.Request().Context(),
		db.InsertTitleParams{
			Title:        titleName,
			NumSubs:      int16(len(stringsSlice)),
			OgLanguageID: fileLangId,
		})
	if err != nil {
		e.Logger().Error(err)
		return nil, nil, err
	}

	// insert phrases into db as translates object of OgLanguage
	_, err = s.translates.InsertNewPhrases(e, title, s.queries, stringsSlice)
	if err != nil {
		e.Logger().Error(err)
		_ = s.queries.DeleteTitleById(e.Request().Context(), title.ID)
		return nil, nil, err
	}
	return &title, nil, nil
}

// AudioFromTitle accepts a title id and returns a zip file of mp3 audio tracks for
// learning a language that you choose
func (s *Server) AudioFromTitle(e echo.Context) error {
	var audioFromTitleRequest oapi.AudioFromTitleJSONRequestBody
	err := e.Bind(&audioFromTitleRequest)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	if audioFromTitleRequest.Pattern != nil {
		s.config.AudioPattern = *audioFromTitleRequest.Pattern
	}
	e.Set("pattern", s.config.AudioPattern)

	if audioFromTitleRequest.Pause != nil {
		s.config.PhrasePause = *audioFromTitleRequest.Pause
	}
	e.Set("pause", audioFromTitleRequest.Pause)

	// get title to translate from
	title, err := s.queries.SelectTitleById(e.Request().Context(), audioFromTitleRequest.TitleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.String(http.StatusBadRequest, "invalid title id")
		}
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	zipFile, err := s.createAudioFromTitle(e, title, audioFromTitleRequest)
	if err != nil {
		if errors.Is(err, util.ErrVoiceLangIdNoMatch) {
			return e.String(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, util.ErrVoiceIdInvalid) {
			return e.String(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, util.ErrOneFile) {
			return e.Attachment(zipFile.Name(), title.Title+".mp3")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.Attachment(zipFile.Name(), title.Title+".zip")
}

// createAudioFromTitle is a helper function that performs the tasks shared by
// AudioFromFile and AudioFromTitle
func (s *Server) createAudioFromTitle(e echo.Context, title db.Title, r oapi.AudioFromTitleJSONRequestBody) (*os.File, error) {
	fromVoice, err := s.queries.SelectVoiceById(e.Request().Context(), r.FromVoiceId)
	if err != nil {
		e.Logger().Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.ErrVoiceIdInvalid
		}
		return nil, err
	}

	toVoice, err := s.queries.SelectVoiceById(e.Request().Context(), r.ToVoiceId)
	if err != nil {
		e.Logger().Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.ErrVoiceIdInvalid
		}
		return nil, err
	}

	// TODO if you don't want these files to persist then you need to defer removing them from calling function
	audioBasePath := fmt.Sprintf("%s%d/", s.config.TTSBasePath, title.ID)

	fromAudioBasePath := fmt.Sprintf("%s%d/", audioBasePath, fromVoice.LanguageID)
	toAudioBasePath := fmt.Sprintf("%s%d/", audioBasePath, toVoice.LanguageID)

	if err = s.translates.CreateTTS(e, s.queries, title, r.FromVoiceId, fromAudioBasePath); err != nil {
		e.Logger().Error(err)
		osErr := os.RemoveAll(audioBasePath)
		if osErr != nil {
			e.Logger().Error(osErr)
		}
		return nil, err
	}

	if err = s.translates.CreateTTS(e, s.queries, title, r.ToVoiceId, toAudioBasePath); err != nil {
		e.Logger().Error(err)
		osErr := os.RemoveAll(audioBasePath)
		if osErr != nil {
			e.Logger().Error(osErr)
		}
		return nil, err
	}

	phraseIds, err := s.queries.SelectPhraseIdsByTitleId(e.Request().Context(), title.ID)
	if err != nil {
		e.Logger().Error(err)
		osErr := os.RemoveAll(audioBasePath)
		if osErr != nil {
			e.Logger().Error(osErr)
		}
		return nil, err
	}

	pausePath, ok := audiofile.AudioPauseFilePath[s.config.PhrasePause]
	if !ok {
		e.Logger().Error(errors.New("audio pause file not found"))
		return nil, err
	}
	fullPausePath := s.config.TTSBasePath + pausePath

	tmpDirPath := fmt.Sprintf("%s%s-%s/", s.config.TTSBasePath, title.Title, test.RandomString(4))
	err = os.MkdirAll(tmpDirPath, 0777)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	if err = s.af.BuildAudioInputFiles(e, phraseIds, title, fullPausePath, fromAudioBasePath, toAudioBasePath, tmpDirPath); err != nil {
		return nil, err
	}

	return s.af.CreateMp3Zip(e, title, tmpDirPath)
}
