package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetLanguages returns a list of all the available languages for the
// translate functions
func (s *Server) GetLanguages(e echo.Context) error {
	languages, err := s.queries.ListLanguages(e.Request().Context())

	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, languages)
}
