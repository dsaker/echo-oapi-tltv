package api

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
	"net/http"
	"os"
	"strconv"
	db "talkliketv.click/tltv/db/sqlc"
)

// FindTitles implements all the handlers in the ServerInterface
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

func (s *Server) AddTitle(eCtx echo.Context) error {

	// get lang id and title from multipart form
	lang := eCtx.FormValue("languageId")
	titleName := eCtx.FormValue("titleName")
	langIdInt16, err := strconv.ParseInt(lang, 10, 16)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	// Get file handler for filename, size and headers
	file, err := eCtx.FormFile("filePath")
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	// Check if file size is too large 32000 == 4KB ~ approximately 2 pages of text
	if file.Size > s.config.FileUploadLimit*8000 {
		return eCtx.String(http.StatusBadRequest, "file too large")
	}
	src, err := file.Open()
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	//Get language model from id
	langModel, err := s.queries.SelectLanguagesById(eCtx.Request().Context(), int16(langIdInt16))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return eCtx.String(http.StatusBadRequest, err.Error())
		}
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	// Create strings slice and count number of lines form titles model
	scanner := bufio.NewScanner(src)
	var stringsSlice []string
	numLines := 0
	for scanner.Scan() {
		numLines += 1
		stringsSlice = append(stringsSlice, scanner.Text())
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	title, err := s.queries.InsertTitle(
		eCtx.Request().Context(),
		db.InsertTitleParams{
			Title:        titleName,
			NumSubs:      int16(numLines),
			OgLanguageID: int16(langIdInt16),
		})

	// use helper function so you can roll back InsertTitle in case of any error
	err = addTitleHelper(eCtx, s, stringsSlice, title, langModel.Tag)
	if err != nil {
		err = s.queries.DeleteTitleById(eCtx.Request().Context(), title.ID)
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusOK, title)
}

func addTitleHelper(eCtx echo.Context, s *Server, slice []string, t db.Title, tag string) error {

	// insert phrases into db as translates object of OgLanguage
	translatesSlice, err := s.translates.InsertNewPhrases(eCtx, t, s.queries, slice)
	if err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	//create base path for storing mp3 audio files
	audioBasePath := s.config.TTSBasePath +
		strconv.Itoa(int(t.ID)) + "/" +
		strconv.Itoa(int(t.OgLanguageID)) + "/"
	// TODO change permission or make permission configurable. Can't test if set to 0644
	err = os.MkdirAll(audioBasePath, 0777)
	if err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	// TODO
	err = s.translates.TextToSpeech(eCtx, translatesSlice, audioBasePath, tag)
	if err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	return nil
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

	// get language tag to translate to
	langTag, err := language.Parse(dbLang.Tag)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
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
	newTranslates, err := s.translates.TranslatePhrases(eCtx, phrasesToTranslate, langTag)
	if err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	// insert new translated phrases into the database
	insertTranslates, err := s.translates.InsertTranslates(eCtx, s.queries, dbLang.ID, newTranslates)
	if err != nil {
		eCtx.Logger().Info(fmt.Sprintf("Error inserting translates -- titleId: %d -- languageId: %d -- error: %s", title.ID, dbLang.ID, err.Error()))
		// TODO		delete translates where language id and phrase id in ( select by title id)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusCreated, insertTranslates)
}
