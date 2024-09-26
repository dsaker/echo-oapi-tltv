package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
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

func (p *Server) CreateUser(ctx echo.Context) error {
	// We expect a NewUser object in the request body.
	var newUser NewUser
	err := ctx.Bind(&newUser)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// We're always asynchronous, so lock unsafe operations below
	p.Lock()
	defer p.Unlock()

	password, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	user, err := p.queries.InsertUser(
		ctx.Request().Context(),
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
				return ctx.String(http.StatusBadRequest, "a user with this email address already exists")
			}
			if db.PqErrorConstraint(err) == db.UsernameConstraint {
				return ctx.String(http.StatusBadRequest, "a user with this name already exists")
			}
			return ctx.String(http.StatusBadRequest, "duplicate key violation")
		}
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	permission, err := p.queries.SelectPermissionByCode(ctx.Request().Context(), db.ReadTitlesCode)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	_, err = p.queries.InsertUserPermission(
		ctx.Request().Context(), db.InsertUserPermissionParams{
			UserID:       user.ID,
			PermissionID: permission.ID,
		})

	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	rsp := newUserResponse(user)
	return ctx.JSON(http.StatusOK, rsp)
}

func (p *Server) DeleteUser(ctx echo.Context, id int64) error {
	err := token.CheckJWTUserIDFromRequest(ctx, id)
	if err != nil {
		return ctx.String(http.StatusForbidden, "Invalid user ID")
	}

	err = p.queries.DeleteUserById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Error deleting user: %s", err))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (p *Server) FindUserByID(ctx echo.Context, id int64) error {
	err := token.CheckJWTUserIDFromRequest(ctx, id)
	if err != nil {
		return ctx.String(http.StatusForbidden, err.Error())
	}

	user, err := p.queries.SelectUserById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Error selecting user by id: %s", err))
	}

	rsp := newUserResponse(user)
	return ctx.JSON(http.StatusOK, rsp)
}

func (p *Server) UpdateUser(ctx echo.Context, id int64) error {

	err := token.CheckJWTUserIDFromRequest(ctx, id)
	if err != nil {
		return ctx.String(http.StatusForbidden, "Invalid user ID")
	}

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	patch, err := jsonpatch.DecodePatch(body)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	user, err := p.queries.SelectUserById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	current := NewUser{
		Email:         user.Email,
		Flipped:       user.Flipped,
		NewLanguageId: user.NewLanguageID,
		OgLanguageId:  user.OgLanguageID,
		Password:      user.HashedPassword,
		TitleId:       user.TitleID,
	}

	currentBytes, err := json.Marshal(current)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	modifiedBytes, err := patch.Apply(currentBytes)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	var modified NewUser
	err = json.Unmarshal(modifiedBytes, &modified)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// perform business logic checks
	if modified.Password != current.Password {
		password, err := bcrypt.GenerateFromPassword([]byte(modified.Password), 14)
		if err != nil {
			ctx.Logger().Error(err)
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		modified.Password = string(password)
	}

	updatedUser, err := p.queries.UpdateUserById(
		ctx.Request().Context(),
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
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	rsp := newUserResponse(updatedUser)
	return ctx.JSON(http.StatusOK, rsp)
}

func (p *Server) LoginUser(ctx echo.Context) error {

	// We expect a NewUser object in the request body.
	var userLogin UserLogin
	err := ctx.Bind(&userLogin)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	user, err := p.queries.SelectUserByName(ctx.Request().Context(), userLogin.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.String(http.StatusNotFound, "username not found")
		}
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	err = util.CheckPassword(userLogin.Password, user.HashedPassword)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "invalid username or password")
	}

	permissions, err := p.queries.SelectUserPermissions(ctx.Request().Context(), user.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctx.Logger().Error(err)
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
	}

	jwsToken, err := p.fa.CreateJWSWithClaims(permissions, user)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, jwsToken)
}
