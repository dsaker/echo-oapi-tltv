package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) GetLanguages(ctx echo.Context) error {
	languages, err := s.Queries.ListLanguages(ctx.Request().Context())

	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, languages)
}
