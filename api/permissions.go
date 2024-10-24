package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
)

func (s *Server) AddUserPermission(ctx echo.Context) error {
	// We expect a NewTitle object in the request body.
	var newUserPermission NewUserPermission
	err := ctx.Bind(&newUserPermission)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	userPermission, err := s.Queries.InsertUserPermission(
		ctx.Request().Context(),
		db.InsertUserPermissionParams{
			UserID:       newUserPermission.UserId,
			PermissionID: newUserPermission.PermissionId,
		})

	if err != nil {
		if db.PqErrorCode(err) == db.ForeignKeyViolation {
			return ctx.String(http.StatusBadRequest, err.Error())
		}
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, userPermission)
}
