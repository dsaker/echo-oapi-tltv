package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/audio/audiofile"
)

// FindTitles returns the number of titles set by the Limit and Similarity params
func (s *Server) FindTitles(ctx echo.Context, params FindTitlesParams) error {

	titles, err := s.queries.ListTitles(
		ctx.Request().Context(),
		db.ListTitlesParams{
			Similarity: params.Similarity,
			Limit:      params.Limit,
		})

	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, titles)
}

// AddTitle takes your uploaded file, filename, and title and adds it to the database,
// along with adding the phrases in the original language and getting the text-to-speech
// and saving it to the file path set in config.TTSBasePath
func (s *Server) AddTitle(e echo.Context) error {

	// get lang id and title from multipart form
	lang := e.FormValue("languageId")
	titleName := e.FormValue("titleName")
	langIdInt16, err := strconv.ParseInt(lang, 10, 16)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	// Get file handler for filename, size and headers
	file, err := e.FormFile("filePath")
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	// Check if file size is too large 32000 == 4KB ~ approximately 2 pages of text
	if file.Size > s.config.FileUploadLimit*8000 {
		return e.String(http.StatusBadRequest, "file too large")
	}
	src, err := file.Open()
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	// Create strings slice and count number of lines form titles model
	stringsSlice, err := audiofile.GetLines(e, src)
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("unable to parse file: %s", err.Error()))
	}
	// TODO add max number of phrases to configs
	if len(stringsSlice) > 800 {
		rString := fmt.Sprintf("file too large, limit is %d, your file has %d lines", 800, len(stringsSlice))
		return e.String(http.StatusBadRequest, rString)
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	title, err := s.queries.InsertTitle(
		e.Request().Context(),
		db.InsertTitleParams{
			Title:        titleName,
			NumSubs:      int16(len(stringsSlice)),
			OgLanguageID: int16(langIdInt16),
		})
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// insert phrases into db as translates object of OgLanguage
	_, err = s.Translates.InsertNewPhrases(e, title, s.queries, stringsSlice)
	if err != nil {
		dbErr := s.queries.DeleteTitleById(e.Request().Context(), title.ID)
		if dbErr != nil {
			return e.String(http.StatusInternalServerError, dbErr.Error())
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, title)
}

func (s *Server) FindTitleByID(ctx echo.Context, id int64) error {

	title, err := s.queries.SelectTitleById(ctx.Request().Context(), id)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, title)
}

func (s *Server) DeleteTitle(ctx echo.Context, id int64) error {

	err := s.queries.DeleteTitleById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (s *Server) TitlesTranslate(eCtx echo.Context) error {

	var newTranslateTitle TitlesTranslateRequest
	err := eCtx.Bind(&newTranslateTitle)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	// make sure the translates for that title don't already exist
	exists, err := s.queries.SelectExistsTranslates(
		eCtx.Request().Context(),
		db.SelectExistsTranslatesParams{
			LanguageID: newTranslateTitle.NewLanguageId,
			ID:         newTranslateTitle.TitleId,
		})
	if exists {
		return eCtx.String(http.StatusBadRequest, "title already exists in that language")
	}

	// get title to translate from
	title, err := s.queries.SelectTitleById(eCtx.Request().Context(), newTranslateTitle.TitleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return eCtx.String(http.StatusBadRequest, "invalid title id")
		}
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	// get language model for tag
	dbLang, err := s.queries.SelectLanguagesById(eCtx.Request().Context(), newTranslateTitle.NewLanguageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return eCtx.String(http.StatusBadRequest, "invalid language id")
		}
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	// get translates for original language to translate from
	phrasesToTranslate, err := s.queries.SelectTranslatesByTitleIdLangId(
		eCtx.Request().Context(),
		db.SelectTranslatesByTitleIdLangIdParams{
			ID:         newTranslateTitle.TitleId,
			LanguageID: title.OgLanguageID,
		})
	if err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	// get the translates from the phrases
	newTranslates, err := s.Translates.TranslatePhrases(eCtx, phrasesToTranslate, dbLang)
	if err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	// insert new translated phrases into the database
	insertTranslates, err := s.Translates.InsertTranslates(eCtx, s.queries, dbLang.ID, newTranslates)
	if err != nil {
		eCtx.Logger().Info(fmt.Sprintf("Error inserting translates -- titleId: %d -- languageId: %d -- error: %s", title.ID, dbLang.ID, err.Error()))
		// TODO		delete translates where language id and phrase id in ( select by title id)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusCreated, insertTranslates)
}
