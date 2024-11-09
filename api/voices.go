package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"talkliketv.click/tltv/internal/oapi"
)

// GetVoices returns a list of all the available voices for the
// text-to-speech functions
func (s *Server) GetVoices(e echo.Context, params oapi.GetVoicesParams) error {

	if params.LanguageId == nil {
		voices, err := s.queries.ListVoices(e.Request().Context())
		if err != nil {
			e.Logger().Error(err)
			return e.String(http.StatusInternalServerError, err.Error())
		}

		return e.JSON(http.StatusOK, voices)
	}

	voices, err := s.queries.SelectVoicesByLanguageId(e.Request().Context(), *params.LanguageId)

	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, voices)
}
