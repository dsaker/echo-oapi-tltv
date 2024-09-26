package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	mockdb "talkliketv.click/tltv/db/mock"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

type eqCreateUserParamsMatcher struct {
	arg      db.InsertUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.InsertUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.InsertUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

type testCase struct {
	name          string
	body          map[string]any
	stringBody    string
	user          db.User
	userId        int64
	buildStubs    func(store *mockdb.MockQuerier)
	checkResponse func(rec *httptest.ResponseRecorder)
}

func TestGetUser(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	testCases := []testCase{
		{
			name:   "get user1",
			user:   user1,
			userId: user1.ID,
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					SelectUserById(gomock.Any(), user1.ID).
					Times(1).
					Return(user1, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				var gotUser db.User
				err := json.Unmarshal([]byte(rec.Body.String()), &gotUser)
				require.NoError(t, err)
				requireMatchAnyExcept(t, user1, gotUser, []string{"HashedPassword", "ID"}, "", "")
			},
		},
		{
			name:   "Id's don't match",
			user:   user1,
			userId: user2.ID,
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, rec.Code)
				require.Contains(t, rec.Body.String(), "provided id does not match user id")
			},
		},
		{
			name:   "User does not exist",
			user:   user2,
			userId: user2.ID,
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					SelectUserById(gomock.Any(), user2.ID).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
				require.Contains(t, rec.Body.String(), "Error selecting user by id: sql: no rows in result set")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			srv, c, rec := setupTest(t, ctrl, tc, string(data), http.MethodGet)

			err = srv.FindUserByID(c, tc.userId)
			require.NoError(t, err)
			tc.checkResponse(rec)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	testCases := []testCase{
		{
			name:   "delete user1",
			user:   user1,
			userId: user1.ID,
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					DeleteUserById(gomock.Any(), user1.ID).
					Times(1).
					Return(nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
		{
			name:   "Id's don't match",
			user:   user1,
			userId: user2.ID,
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, rec.Code)
				require.Contains(t, rec.Body.String(), "Invalid user ID")
			},
		},
		{
			name:   "User does not exist",
			user:   user2,
			userId: user2.ID,
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					DeleteUserById(gomock.Any(), user2.ID).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
				require.Contains(t, rec.Body.String(), "Error deleting user: sql: no rows in result set")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			srv, c, rec := setupTest(t, ctrl, tc, string(data), http.MethodDelete)

			err = srv.DeleteUser(c, tc.userId)
			require.NoError(t, err)
			tc.checkResponse(rec)
		})
	}
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)
	testCases := []testCase{
		{
			name: "Create User",
			body: map[string]any{
				"email":         user.Email,
				"flipped":       user.Flipped,
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      password,
				"titleId":       user.TitleID,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				arg := db.InsertUserParams{
					Name:           user.Name,
					Email:          user.Email,
					HashedPassword: password,
					TitleID:        user.TitleID,
					Flipped:        user.Flipped,
					OgLanguageID:   user.OgLanguageID,
					NewLanguageID:  user.NewLanguageID,
				}
				arg2 := db.InsertUserPermissionParams{
					UserID:       user.ID,
					PermissionID: util.ValidPermissionId,
				}
				store.EXPECT().
					InsertUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Return(user, nil)
				store.EXPECT().
					SelectPermissionByCode(
						gomock.Any(),
						db.ReadTitlesCode).
					Return(db.Permission{ID: util.ValidPermissionId, Code: ""}, nil)
				store.EXPECT().
					InsertUserPermission(gomock.Any(), arg2).
					Times(1).
					Return(db.UsersPermission{UserID: user.ID, PermissionID: util.ValidPermissionId}, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)

				var gotUser db.User
				err := json.Unmarshal([]byte(rec.Body.String()), &gotUser)
				require.NoError(t, err)
				requireMatchAnyExcept(t, user, gotUser, []string{"HashedPassword", "ID"}, "", "")
			},
		},
		{
			name: "Internal Server",
			body: map[string]any{
				"email":         user.Email,
				"flipped":       user.Flipped,
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      password,
				"titleId":       user.TitleID,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					InsertUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
				require.Contains(t, rec.Body.String(), "sql: connection is already closed")
			},
		},
		{
			name: "Duplicate Username",
			body: map[string]any{
				"email":         user.Email,
				"flipped":       user.Flipped,
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      password,
				"titleId":       user.TitleID,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					InsertUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrUniqueViolation)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
				require.Contains(t, rec.Body.String(), "duplicate key violation")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			srv, c, rec := setupTest(t, ctrl, tc, string(data), http.MethodPut)

			err = srv.CreateUser(c)
			require.NoError(t, err)
			tc.checkResponse(rec)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	user1, _ := randomUser(t)
	updateUserParams := db.UpdateUserByIdParams{
		Email:          user1.Email,
		TitleID:        user1.TitleID,
		Flipped:        user1.Flipped,
		OgLanguageID:   user1.OgLanguageID,
		NewLanguageID:  user1.NewLanguageID,
		HashedPassword: user1.HashedPassword,
		ID:             user1.ID,
	}

	user2, _ := randomUser(t)

	testCases := []testCase{
		{
			name:   "update user1 email",
			user:   user1,
			userId: user1.ID,
			stringBody: `[
			{
				"op": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`,
			buildStubs: func(store *mockdb.MockQuerier) {
				userCopy := user1
				paramsCopy := updateUserParams
				paramsCopy.Email = "newemail2@email.com"
				userCopy.Email = "newemail2@email.com"
				store.EXPECT().
					SelectUserById(gomock.Any(), user1.ID).
					Times(1).
					Return(user1, nil)
				store.EXPECT().
					UpdateUserById(gomock.Any(), paramsCopy).
					Times(1).
					Return(userCopy, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				var gotUser db.User
				err := json.Unmarshal([]byte(rec.Body.String()), &gotUser)
				require.NoError(t, err)
				requireMatchAnyExcept(t, user1, gotUser, []string{"HashedPassword", "ID"}, "Email", "newemail2@email.com")
			},
		},
		{
			name:   "Id does not match",
			user:   user1,
			userId: user2.ID,
			stringBody: `[
			{
				"op": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`,
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, rec.Code)
				require.Contains(t, rec.Body.String(), "Invalid user ID")
			},
		},
		{
			name:   "Invalid format for PatchUser",
			user:   user1,
			userId: user1.ID,
			stringBody: `[
			{
				"wrong": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`,
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
				require.Contains(t, rec.Body.String(), "invalid operation {\"path\":\"/email\",\"value\":\"newemail2@email.com\",\"wrong\":\"replace\"}: unsupported operation")
			},
		},
		{
			name:   "User does not exist",
			user:   user2,
			userId: user2.ID,
			stringBody: `[
			{
				"op": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`,
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					SelectUserById(gomock.Any(), user2.ID).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
				require.Contains(t, rec.Body.String(), "Error selecting user by id: sql: no rows in result set")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv, c, rec := setupTest(t, ctrl, tc, tc.stringBody, http.MethodPut)
			err := srv.UpdateUser(c, tc.userId)
			require.NoError(t, err)
			tc.checkResponse(rec)
		})
	}
}

func TestLoginUser(t *testing.T) {
	user, password := randomUser(t)
	permissions := []string{db.ReadTitlesCode}

	testCases := []testCase{
		{
			name: "OK",
			body: map[string]any{
				"username": user.Name,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					SelectUserByName(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					SelectUserPermissions(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(permissions, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "UserNotFound",
			body: map[string]any{
				"username": "NotFound",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					SelectUserByName(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
				require.Contains(t, rec.Body.String(), "invalid username or password")
			},
		},
		{
			name: "IncorrectPassword",
			body: map[string]any{
				"username": user.Name,
				"password": "incorrect",
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					SelectUserByName(gomock.Any(), gomock.Eq(user.Name)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
				require.Contains(t, rec.Body.String(), "invalid username or password")
			},
		},
		{
			name: "InternalError",
			body: map[string]any{
				"username": user.Name,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					SelectUserByName(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
				require.Contains(t, rec.Body.String(), "sql: connection is already closed")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			srv, c, rec := setupTest(t, ctrl, tc, string(data), http.MethodGet)

			err = srv.LoginUser(c)
			require.NoError(t, err)
			tc.checkResponse(rec)
		})
	}
}

func TestCreateUserMiddleware(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		body          map[string]any
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(res *http.Response)
	}{
		{
			name: "Invalid Username",
			body: map[string]any{
				"name":          "invalid-user#1",
				"email":         user.Email,
				"flipped":       user.Flipped,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      password,
				"titleId":       user.TitleID,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "InvalidEmail",
			body: map[string]any{
				"flipped":       user.Flipped,
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      password,
				"titleId":       user.TitleID,
				"email":         "invalid-email",
			},
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "string doesn't match the regular expression ")
			},
		},
		{
			name: "TooShortPassword",
			body: map[string]any{
				"flipped":       user.Flipped,
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      "string",
				"titleId":       user.TitleID,
				"email":         user.Email,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "minimum string length is 8")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockQuerier(ctrl)
			tc.buildStubs(store)

			spec, err := GetSwagger()
			require.NoError(t, err)

			e := echo.New()
			_ = NewServer(e, testCfg, store, spec)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			ts := httptest.NewServer(e)
			urlPath := "/v1/users"
			res, err := ts.Client().Post(ts.URL+urlPath, "application/json", bytes.NewReader(data))
			require.NoError(t, err)
			tc.checkResponse(res)
		})
	}
}
