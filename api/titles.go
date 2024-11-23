package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/util"
)

// FindTitles returns the number of titles set by the Limit and Similarity params
func (s *Server) FindTitles(ctx echo.Context, params oapi.FindTitlesParams) error {
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
// along with adding the phrases in the original language to the translates table
func (s *Server) AddTitle(e echo.Context) error {
	// get lang id and title from multipart form
	titleName := e.FormValue("titleName")
	fileLangId, err := util.ConvertStringInt16(e.FormValue("languageId"))
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("error converting fileLanguageId to int16: %s", err.Error()))
	}

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

	return e.JSON(http.StatusOK, title)
}

func (s *Server) FindTitleByID(e echo.Context, id int64) error {
	title, err := s.queries.SelectTitleById(e.Request().Context(), id)
	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, title)
}

func (s *Server) DeleteTitle(e echo.Context, id int64) error {
	err := s.queries.DeleteTitleById(e.Request().Context(), id)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}
	return e.NoContent(http.StatusNoContent)
}

// TitlesTranslate translates the phrases of a title from the original language
// of the title to any available language by id and stores them in the translates table
func (s *Server) TitlesTranslate(e echo.Context) error {
	var newTranslateTitle oapi.TitlesTranslateRequest
	err := e.Bind(&newTranslateTitle)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	// make sure the translates for that title don't already exist
	exists, err := s.queries.SelectExistsTranslates(
		e.Request().Context(),
		db.SelectExistsTranslatesParams{
			LanguageID: newTranslateTitle.NewLanguageId,
			ID:         newTranslateTitle.TitleId,
		})
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	if exists {
		return e.String(http.StatusBadRequest, "title already exists in that language")
	}

	// get title to translate from
	title, err := s.queries.SelectTitleById(e.Request().Context(), newTranslateTitle.TitleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.String(http.StatusBadRequest, "invalid title id")
		}
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// get language model for tag
	dbLang, err := s.queries.SelectLanguagesById(e.Request().Context(), newTranslateTitle.NewLanguageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.String(http.StatusBadRequest, "invalid language id")
		}
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	// get translates for original language to translate from
	phrasesToTranslate, err := s.queries.SelectTranslatesByTitleIdLangId(
		e.Request().Context(),
		db.SelectTranslatesByTitleIdLangIdParams{
			ID:         newTranslateTitle.TitleId,
			LanguageID: title.OgLanguageID,
		})
	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// get the translates from the phrases
	newTranslates, err := s.translates.TranslatePhrases(e, phrasesToTranslate, dbLang)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// check if returned translates is empty
	if len(newTranslates) == 0 {
		return e.String(http.StatusInternalServerError, "something went wrong")
	}

	// insert new translated phrases into the database
	insertTranslates, err := s.translates.InsertTranslates(e, s.queries, dbLang.ID, newTranslates)
	if err != nil {
		e.Logger().Info(fmt.Sprintf("Error inserting translates -- titleId: %d -- languageId: %d -- error: %s", title.ID, dbLang.ID, err.Error()))
		// roll back by deleting any translates that were inserted
		_ = s.queries.DeleteTranslatesByLanguageId(e.Request().Context(),
			db.DeleteTranslatesByLanguageIdParams{
				LanguageID: dbLang.ID,
				TitleID:    title.ID,
			})
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, insertTranslates)
}
