package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) AudioFromFile(ctx echo.Context) error {
	ctx.Logger().Info("inside AudioFromFile")
	return ctx.NoContent(http.StatusNotImplemented)
}

func (s *Server) AudioFromTitle(ctx echo.Context) error {
	ctx.Logger().Info("inside AudioFromTitle")
	return ctx.NoContent(http.StatusNotImplemented)
}
