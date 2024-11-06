package api

import (
	"encoding/json"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/token"
)

// GetPhrases gets the next set limit or default ten phrases for user title and language ids
// sorted by ascending correctness (correctness is how many times a user has guessed a phrase
// correctly)
func (s *Server) GetPhrases(e echo.Context, params oapi.GetPhrasesParams) error {

	if params.Limit == nil {
		params.Limit = new(int32)
		*params.Limit = 10
	}

	user, err := token.GetUserFromContext(e)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	usersPhrases, err := s.queries.SelectTranslatesWithCorrect(
		e.Request().Context(),
		db.SelectTranslatesWithCorrectParams{
			UserID:       user.ID,
			TitleID:      user.TitleID,
			LanguageID:   user.OgLanguageID,
			LanguageID_2: user.NewLanguageID,
			Limit:        *params.Limit,
		})

	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, usersPhrases)
}

// UpdateUsersPhrases performs a PATCH request on the users_phrases table. It will
// mostly be used to increase the correct column by one.
func (s *Server) UpdateUsersPhrases(e echo.Context, phraseId int64, languageId int16) error {
	user, err := token.GetUserFromContext(e)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	body, err := io.ReadAll(e.Request().Body)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	patch, err := jsonpatch.DecodePatch(body)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	usersPhraseById, err := s.queries.SelectUsersPhrasesByIds(
		e.Request().Context(),
		db.SelectUsersPhrasesByIdsParams{
			UserID:     user.ID,
			LanguageID: languageId,
			PhraseID:   phraseId,
		})
	if err != nil {
		return e.String(http.StatusInternalServerError, fmt.Sprintf("Error selecting user phrase by ids: %s", err.Error()))
	}

	current := oapi.UsersPhrases{
		LanguageId:    usersPhraseById.LanguageID,
		TitleId:       usersPhraseById.TitleID,
		PhraseId:      usersPhraseById.PhraseID,
		UserId:        usersPhraseById.UserID,
		PhraseCorrect: usersPhraseById.PhraseCorrect,
	}

	currentBytes, err := json.Marshal(current)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	modifiedBytes, err := patch.Apply(currentBytes)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	var modified oapi.UsersPhrases
	err = json.Unmarshal(modifiedBytes, &modified)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	updatedUsersPhrases, err := s.queries.UpdateUsersPhrasesByThreeIds(
		e.Request().Context(),
		db.UpdateUsersPhrasesByThreeIdsParams{
			TitleID:       modified.TitleId,
			UserID:        modified.UserId,
			PhraseID:      modified.PhraseId,
			LanguageID:    modified.LanguageId,
			PhraseCorrect: modified.PhraseCorrect,
		})

	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, updatedUsersPhrases)
}
