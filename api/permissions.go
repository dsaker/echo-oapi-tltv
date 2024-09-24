package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
)

func (p *Server) AddUserPermission(ctx echo.Context) error {
	// We expect a NewTitle object in the request body.
	var newUserPermission NewUserPermission
	err := ctx.Bind(&newUserPermission)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// We're always asynchronous, so lock unsafe operations below
	p.Lock()
	defer p.Unlock()

	userPermission, err := p.queries.InsertUserPermission(
		ctx.Request().Context(),
		db.InsertUserPermissionParams{
			UserID:       newUserPermission.UserId,
			PermissionID: newUserPermission.PermissionId,
		})

	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusCreated, userPermission)
}
