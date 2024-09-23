package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/oapi"
)

func (p *Api) AddUserPermission(w http.ResponseWriter, r *http.Request) {
	// We expect a NewTitle object in the request body.
	var newUserPermission oapi.NewUserPermission
	if err := json.NewDecoder(r.Body).Decode(&newUserPermission); err != nil {
		sendApiError(w, http.StatusBadRequest, "Invalid format for NewUserPermission")
		return
	}

	// We're always asynchronous, so lock unsafe operations below
	p.Lock.Lock()
	defer p.Lock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), p.config.CtxTimeout)
	defer cancel()

	userPermission, err := p.queries.InsertUserPermission(
		ctx, db.InsertUserPermissionParams{
			UserID:       newUserPermission.UserId,
			PermissionID: newUserPermission.PermissionId,
		})

	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error creating user permission: %s", err))
		return
	}

	// Now, we have to return the NewUserPermission
	if err = json.NewEncoder(w).Encode(userPermission); err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error encoding user permission: %s", err))
		return
	}
}
