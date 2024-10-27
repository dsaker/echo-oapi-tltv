package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/oapi"
)

func (s *Server) AddUserPermission(e echo.Context) error {
	// We expect a NewTitle object in the request body.
	var newUserPermission oapi.NewUserPermission
	err := e.Bind(&newUserPermission)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	userPermission, err := s.queries.InsertUserPermission(
		e.Request().Context(),
		db.InsertUserPermissionParams{
			UserID:       newUserPermission.UserId,
			PermissionID: newUserPermission.PermissionId,
		})

	if err != nil {
		if db.PqErrorCode(err) == db.ForeignKeyViolation {
			return e.String(http.StatusBadRequest, err.Error())
		}
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, userPermission)
}
