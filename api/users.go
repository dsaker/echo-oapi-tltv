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
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/token"
	"talkliketv.click/tltv/internal/util"
	"time"
)

type userResponse struct {
	TitleID       int64     `json:"title_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Flipped       bool      `json:"flipped"`
	OgLanguageID  int16     `json:"og_language_id"`
	NewLanguageID int16     `json:"new_language_id"`
	Created       time.Time `json:"created"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		TitleID:       user.TitleID,
		Name:          user.Name,
		Email:         user.Email,
		OgLanguageID:  user.OgLanguageID,
		NewLanguageID: user.NewLanguageID,
		Created:       user.Created,
	}
}

// CreateUser registers a new user
func (s *Server) CreateUser(e echo.Context) error {
	// We expect a NewUser object in the request body.
	var newUser oapi.NewUser
	err := e.Bind(&newUser)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	password, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	user, err := s.queries.InsertUser(
		e.Request().Context(),
		db.InsertUserParams{
			Name:           newUser.Name,
			Email:          newUser.Email,
			HashedPassword: string(password),
			TitleID:        newUser.TitleId,
			OgLanguageID:   newUser.OgLanguageId,
			NewLanguageID:  newUser.NewLanguageId,
		})

	if err != nil {
		if db.PqErrorCode(err) == db.UniqueViolation {
			if db.PqErrorConstraint(err) == db.EmailConstraint {
				return e.String(http.StatusBadRequest, "a user with this email address already exists")
			}
			if db.PqErrorConstraint(err) == db.UsernameConstraint {
				return e.String(http.StatusBadRequest, "a user with this name already exists")
			}
			return e.String(http.StatusBadRequest, "duplicate key violation")
		}
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	rsp := newUserResponse(user)
	return e.JSON(http.StatusOK, rsp)
}

func (s *Server) DeleteUser(e echo.Context, id int64) error {
	err := token.CheckJWTUserIDFromRequest(e, id)
	if err != nil {
		return e.String(http.StatusForbidden, "Invalid user ID")
	}

	err = s.queries.DeleteUserById(e.Request().Context(), id)
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("Error deleting user: %s", err))
	}

	return e.NoContent(http.StatusNoContent)
}

func (s *Server) FindUserByID(e echo.Context, id int64) error {
	err := token.CheckJWTUserIDFromRequest(e, id)
	if err != nil {
		return e.String(http.StatusForbidden, err.Error())
	}

	user, err := s.queries.SelectUserById(e.Request().Context(), id)
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("Error selecting user by id: %s", err))
	}

	rsp := newUserResponse(user)
	return e.JSON(http.StatusOK, rsp)
}

// UpdateUser accepts a Patch request to update the user values
func (s *Server) UpdateUser(e echo.Context, id int64) error {

	err := token.CheckJWTUserIDFromRequest(e, id)
	if err != nil {
		return e.String(http.StatusForbidden, "Invalid user ID")
	}

	body, err := io.ReadAll(e.Request().Body)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	patch, err := jsonpatch.DecodePatch(body)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	user, err := s.queries.SelectUserById(e.Request().Context(), id)
	if err != nil {
		return e.String(http.StatusBadRequest, fmt.Sprintf("Error selecting user by id: %s", err))
	}

	current := oapi.NewUser{
		Email:         user.Email,
		NewLanguageId: user.NewLanguageID,
		OgLanguageId:  user.OgLanguageID,
		Password:      user.HashedPassword,
		TitleId:       user.TitleID,
	}

	currentBytes, err := json.Marshal(current)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	modifiedBytes, err := patch.Apply(currentBytes)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	var modified oapi.NewUser
	err = json.Unmarshal(modifiedBytes, &modified)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	// perform business logic checks
	if modified.Password != current.Password {
		password, err := bcrypt.GenerateFromPassword([]byte(modified.Password), 14)
		if err != nil {
			e.Logger().Error(err)
			return e.String(http.StatusInternalServerError, err.Error())
		}
		modified.Password = string(password)
	}

	updatedUser, err := s.queries.UpdateUserById(
		e.Request().Context(),
		db.UpdateUserByIdParams{
			TitleID:        modified.TitleId,
			Email:          modified.Email,
			OgLanguageID:   modified.OgLanguageId,
			NewLanguageID:  modified.NewLanguageId,
			HashedPassword: modified.Password,
			ID:             id,
		})

	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	rsp := newUserResponse(updatedUser)
	return e.JSON(http.StatusOK, rsp)
}

func (s *Server) LoginUser(e echo.Context) error {

	// We expect a NewUser object in the request body.
	var userLogin oapi.UserLogin
	err := e.Bind(&userLogin)
	if err != nil {
		return e.String(http.StatusBadRequest, err.Error())
	}

	user, err := s.queries.SelectUserByName(e.Request().Context(), userLogin.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.String(http.StatusUnauthorized, "invalid username or password")
		}
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	err = util.CheckPassword(userLogin.Password, user.HashedPassword)
	if err != nil {
		return e.String(http.StatusUnauthorized, "invalid username or password")
	}

	permissions, err := s.queries.SelectUserPermissions(e.Request().Context(), user.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			e.Logger().Error(err)
			return e.String(http.StatusInternalServerError, err.Error())
		}
	}

	jwsToken, err := s.fa.CreateJWSWithClaims(permissions, user)
	if err != nil {
		e.Logger().Error(err)
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, jwsToken)
}
