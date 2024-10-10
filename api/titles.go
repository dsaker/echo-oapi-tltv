package api

import (
	"bufio"
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

	lang := eCtx.FormValue("languageId")
	titleName := eCtx.FormValue("titleName")
	langIdInt16, err := strconv.ParseInt(lang, 10, 16)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	//Get language model from id for tag
	langModel, err := s.queries.SelectLanguagesById(eCtx.Request().Context(), int16(langIdInt16))
	if err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	// Get handler for filename, size and headers
	file, err := eCtx.FormFile("filePath")
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	// Check if file size is too large 32000 == 4KB ~ approximately 4 pages
	if file.Size > 32000 {
		return eCtx.String(http.StatusBadRequest, "file too large")
	}
	src, err := file.Open()
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	eCtx.Logger().Info(fmt.Sprintf("File uploaded successfully: %s", file.Filename))

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

	// insert phrases into db as translates object of OgLanguage
	translatesSlice, err := s.translates.InsertPhrases(eCtx, title, s.queries, stringsSlice, numLines)
	if err != nil {
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	//create base path for storing mp3 audio files
	audioBasePath := s.config.TTSBasePath +
		strconv.Itoa(int(title.ID)) + "/" +
		strconv.Itoa(int(title.OgLanguageID)) + "/"
	err = os.MkdirAll(audioBasePath, os.ModePerm)
	if err != nil {
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}
	err = s.translates.TextToSpeech(eCtx, translatesSlice, audioBasePath, langModel.Tag)

	if err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusOK, title)
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

func (s *Server) TranslateTitle(eCtx echo.Context) error {

	var newTranslateTitle TranslateTitleJSONBody
	err := eCtx.Bind(&newTranslateTitle)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	exists, err := s.queries.SelectExistsTranslates(
		eCtx.Request().Context(),
		db.SelectExistsTranslatesParams{
			LanguageID: newTranslateTitle.NewLanguageId,
			ID:         newTranslateTitle.TitleId,
		})
	if exists {
		return eCtx.String(http.StatusBadRequest, "title already exists in that language")
	}

	title, err := s.queries.SelectTitleById(eCtx.Request().Context(), newTranslateTitle.TitleId)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, "invalid title id")
	}

	dbLang, err := s.queries.SelectLanguagesById(eCtx.Request().Context(), newTranslateTitle.NewLanguageId)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	langTag, err := language.Parse(dbLang.Tag)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	translatesRow, err := s.queries.SelectTranslatesByTitleIdLangId(
		eCtx.Request().Context(),
		db.SelectTranslatesByTitleIdLangIdParams{
			ID:         newTranslateTitle.TitleId,
			LanguageID: title.OgLanguageID,
		})
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	newTranslates, err := s.translates.TranslatePhrases(eCtx, int(translatesRow[0].NumSubs), translatesRow, langTag)
	if err != nil {
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	insertTranslates, err := s.translates.InsertPhrases(eCtx, title, s.queries, newTranslates, int(translatesRow[0].NumSubs))
	if err != nil {
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusCreated, insertTranslates[0])
}
