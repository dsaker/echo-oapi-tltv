package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
	oc "talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/token"
	"talkliketv.click/tltv/internal/util"
	"time"
)

type userResponse struct {
	TitleID       int64     `json:"title_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Flipped       bool      `json:"flipped"`
	OgLanguageID  int64     `json:"og_language_id"`
	NewLanguageID int64     `json:"new_language_id"`
	Created       time.Time `json:"created"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		TitleID:       user.TitleID,
		Name:          user.Name,
		Email:         user.Email,
		Flipped:       user.Flipped,
		OgLanguageID:  user.OgLanguageID,
		NewLanguageID: user.NewLanguageID,
		Created:       user.Created,
	}
}

func (p *Api) Register(w http.ResponseWriter, r *http.Request) {
	// We expect a NewUser object in the request body.
	var newUser oc.NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		sendApiError(w, http.StatusBadRequest, "Invalid format for NewUser")
		return
	}

	// We're always asynchronous, so lock unsafe operations below
	p.Lock.Lock()
	defer p.Lock.Unlock()

	password, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		sendApiError(w, http.StatusBadRequest, "Error generating password")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.CtxTimeout)
	defer cancel()

	user, err := p.queries.InsertUser(
		ctx,
		db.InsertUserParams{
			Name:           newUser.Name,
			Email:          newUser.Email,
			HashedPassword: string(password),
			TitleID:        newUser.TitleId,
			Flipped:        newUser.Flipped,
			OgLanguageID:   newUser.OgLanguageId,
			NewLanguageID:  newUser.NewLanguageId,
		})

	if err != nil {
		if db.PqErrorCode(err) == db.UniqueViolation {
			if db.PqErrorConstraint(err) == db.EmailConstraint {
				sendApiError(w, http.StatusBadRequest, "a user with this email address already exists")
				return
			}
			if db.PqErrorConstraint(err) == db.UsernameConstraint {
				sendApiError(w, http.StatusBadRequest, "a user with this name already exists")
				return
			}
			sendApiError(w, http.StatusBadRequest, "duplicate key violation")
			return
		}
		p.sendInternalError(w, err)
		return
	}

	permission, err := p.queries.SelectPermissionByCode(ctx, db.ReadTitlesCode)
	if err != nil {
		p.sendInternalError(w, err)
		return
	}
	_, err = p.queries.InsertUserPermission(
		ctx, db.InsertUserPermissionParams{
			UserID:       user.ID,
			PermissionID: permission.ID,
		})

	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error creating user permission: %s", err))
		return
	}

	rsp := newUserResponse(user)
	if err = json.NewEncoder(w).Encode(rsp); err != nil {
		sendApiError(w, http.StatusBadRequest, "Error encoding user")
		return
	}
}

func (p *Api) DeleteUser(w http.ResponseWriter, r *http.Request, id int64) {
	err := token.CheckJWTUserIDFromRequest(r.Context(), id)
	if err != nil {
		sendApiError(w, http.StatusForbidden, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.CtxTimeout)
	defer cancel()

	err = p.queries.DeleteUserById(ctx, id)
	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error deleting user: %s", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (p *Api) FindUserByID(w http.ResponseWriter, r *http.Request, id int64) {
	err := token.CheckJWTUserIDFromRequest(r.Context(), id)
	if err != nil {
		sendApiError(w, http.StatusForbidden, "Invalid user ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.CtxTimeout)
	defer cancel()

	user, err := p.queries.SelectUserById(ctx, id)
	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error selecting user by id: %s", err))
		return
	}

	rsp := newUserResponse(user)
	if err = json.NewEncoder(w).Encode(rsp); err != nil {
		sendApiError(w, http.StatusBadRequest, "Error encoding user")
		return
	}
}

func (p *Api) UpdateUser(w http.ResponseWriter, r *http.Request, id int64) {

	err := token.CheckJWTUserIDFromRequest(r.Context(), id)
	if err != nil {
		sendApiError(w, http.StatusForbidden, "Invalid user ID")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendApiError(w, http.StatusBadRequest, "Invalid format for PatchUser")
		return
	}

	patch, err := jsonpatch.DecodePatch(body)
	if err != nil {
		sendApiError(w, http.StatusBadRequest, "Invalid format for PatchUser")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.CtxTimeout)
	defer cancel()

	user, err := p.queries.SelectUserById(ctx, id)
	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error selecting user by id: %s", err))
		return
	}

	current := oc.NewUser{
		Email:         user.Email,
		Flipped:       user.Flipped,
		NewLanguageId: user.NewLanguageID,
		OgLanguageId:  user.OgLanguageID,
		Password:      user.HashedPassword,
		TitleId:       user.TitleID,
	}

	currentBytes, err := json.Marshal(current)
	if err != nil {
		p.sendInternalError(w, err)
		return
	}

	modifiedBytes, err := patch.Apply(currentBytes)
	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error patching user: %s", err))
		return
	}

	var modified oc.NewUser
	err = json.Unmarshal(modifiedBytes, &modified)
	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error unmarshalling modified user: %s", err))
		return
	}

	// perform business logic checks
	if modified.Password != current.Password {
		password, err := bcrypt.GenerateFromPassword([]byte(modified.Password), 14)
		if err != nil {
			p.sendInternalError(w, err)
			return
		}
		modified.Password = string(password)
	}

	updatedUser, err := p.queries.UpdateUserById(
		ctx,
		db.UpdateUserByIdParams{
			TitleID:        modified.TitleId,
			Email:          modified.Email,
			Flipped:        modified.Flipped,
			OgLanguageID:   modified.OgLanguageId,
			NewLanguageID:  modified.NewLanguageId,
			HashedPassword: modified.Password,
			ID:             id,
		})

	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error updating user: %s", err))
		return
	}

	rsp := newUserResponse(updatedUser)
	if err = json.NewEncoder(w).Encode(rsp); err != nil {
		sendApiError(w, http.StatusBadRequest, "Error encoding user")
		return
	}
}

func (p *Api) LoginUser(w http.ResponseWriter, r *http.Request) {

	// We expect a NewUser object in the request body.
	var userLogin oc.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		sendApiError(w, http.StatusBadRequest, "Invalid format for UserLogin")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.CtxTimeout)
	defer cancel()

	user, err := p.queries.SelectUserByName(ctx, userLogin.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendApiError(w, http.StatusNotFound, "Invalid username or password")
			return
		}
		p.sendInternalError(w, err)
		return
	}

	err = util.CheckPassword(userLogin.Password, user.HashedPassword)
	if err != nil {
		sendApiError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	permissions, err := p.queries.SelectUserPermissions(ctx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendApiError(w, http.StatusNotFound, err.Error())
			return
		}
		p.logger.PrintError(err, nil)
		sendApiError(w, http.StatusInternalServerError, "Error selecting user")
		return
	}

	jwsToken, err := p.fa.CreateJWSWithClaims(permissions, user)
	if err != nil {
		p.logger.PrintError(err, nil)
		sendApiError(w, http.StatusInternalServerError, "Error creating jwsToken")
		return
	}

	if err = json.NewEncoder(w).Encode(jwsToken); err != nil {
		sendApiError(w, http.StatusBadRequest, "Error encoding user")
		return
	}
}
