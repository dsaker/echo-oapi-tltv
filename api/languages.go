package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"talkliketv.click/tltv/internal/oapi"
)

// GetLanguages returns a list of all the available languages for the
// translate functions
func (s *Server) GetLanguages(e echo.Context, params oapi.GetLanguagesParams) error {

	similarity := ""
	if params.Similarity != nil {
		similarity = *params.Similarity
	}
	languages, err := s.queries.ListLanguagesSimilar(e.Request().Context(), similarity)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, languages)
}
