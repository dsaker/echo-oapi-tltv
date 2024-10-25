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

func (s *Server) GetPhrases(ctx echo.Context, params oapi.GetPhrasesParams) error {

	if params.Limit == nil {
		params.Limit = new(int32)
		*params.Limit = 10
	}

	user, err := token.GetUserFromContext(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	usersPhrases, err := s.queries.SelectPhrasesFromTranslatesWithCorrect(
		ctx.Request().Context(),
		db.SelectPhrasesFromTranslatesWithCorrectParams{
			UserID:       user.ID,
			TitleID:      user.TitleID,
			LanguageID:   user.OgLanguageID,
			LanguageID_2: user.NewLanguageID,
			Limit:        *params.Limit,
		})

	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, usersPhrases)
}

func (s *Server) UpdateUsersPhrases(ctx echo.Context, phraseId int64, languageId int16) error {
	user, err := token.GetUserFromContext(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	patch, err := jsonpatch.DecodePatch(body)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	usersPhraseById, err := s.queries.SelectUsersPhrasesByIds(
		ctx.Request().Context(),
		db.SelectUsersPhrasesByIdsParams{
			UserID:     user.ID,
			LanguageID: languageId,
			PhraseID:   phraseId,
		})
	if err != nil {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error selecting user phrase by ids: %s", err.Error()))
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
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	modifiedBytes, err := patch.Apply(currentBytes)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	var modified oapi.UsersPhrases
	err = json.Unmarshal(modifiedBytes, &modified)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	updatedUsersPhrases, err := s.queries.UpdateUsersPhrasesByThreeIds(
		ctx.Request().Context(),
		db.UpdateUsersPhrasesByThreeIdsParams{
			TitleID:       modified.TitleId,
			UserID:        modified.UserId,
			PhraseID:      modified.PhraseId,
			LanguageID:    modified.LanguageId,
			PhraseCorrect: modified.PhraseCorrect,
		})

	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, updatedUsersPhrases)
}
