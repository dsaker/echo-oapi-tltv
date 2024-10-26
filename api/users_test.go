package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	mockdb "talkliketv.click/tltv/db/mock"
	db "talkliketv.click/tltv/db/sqlc"
	mockt "talkliketv.click/tltv/internal/mock/translates"
	"talkliketv.click/tltv/internal/test"
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

func TestGetUser(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	testCases := []testCase{
		{
			name:   "get user1",
			user:   user1,
			userId: user1.ID,
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				store.EXPECT().
					SelectUserById(gomock.Any(), user1.ID).
					Times(1).
					Return(user1, nil)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				var gotUser db.User
				err := json.Unmarshal([]byte(rec.Body.String()), &gotUser)
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, user1, gotUser, []string{"HashedPassword", "ID"}, "", "")
			},
		},
		{
			name:   "Id's don't match",
			user:   user1,
			userId: user2.ID,
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, rec.Code)
				require.Contains(t, rec.Body.String(), "provided id does not match user id")
			},
		},
		{
			name:   "User does not exist",
			user:   user2,
			userId: user2.ID,
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				store.EXPECT().
					SelectUserById(gomock.Any(), user2.ID).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
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

			urlPath := usersBasePath + "/" + strconv.FormatInt(tc.user.ID, 10)
			srv, c, rec := setupHandlerTest(t, ctrl, tc, urlPath, string(data), http.MethodGet)

			err = srv.FindUserByID(c, tc.userId)
			require.NoError(t, err)
			tc.checkRecorder(rec)
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
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				store.EXPECT().
					DeleteUserById(gomock.Any(), user1.ID).
					Times(1).
					Return(nil)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, rec.Code)
			},
		},
		{
			name:   "Id's don't match",
			user:   user1,
			userId: user2.ID,
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, rec.Code)
				require.Contains(t, rec.Body.String(), "Invalid user ID")
			},
		},
		{
			name:   "User does not exist",
			user:   user2,
			userId: user2.ID,
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				store.EXPECT().
					DeleteUserById(gomock.Any(), user2.ID).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
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

			urlPath := usersBasePath + "/" + strconv.FormatInt(tc.user.ID, 10)
			srv, c, rec := setupHandlerTest(t, ctrl, tc, urlPath, string(data), http.MethodDelete)

			err = srv.DeleteUser(c, tc.userId)
			require.NoError(t, err)
			tc.checkRecorder(rec)
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
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      password,
				"titleId":       user.TitleID,
			},
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				arg := db.InsertUserParams{
					Name:           user.Name,
					Email:          user.Email,
					HashedPassword: password,
					TitleID:        user.TitleID,
					OgLanguageID:   user.OgLanguageID,
					NewLanguageID:  user.NewLanguageID,
				}
				arg2 := db.InsertUserPermissionParams{
					UserID:       user.ID,
					PermissionID: test.ValidPermissionId,
				}
				store.EXPECT().
					InsertUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Return(user, nil)
				store.EXPECT().
					SelectPermissionByCode(
						gomock.Any(),
						db.ReadTitlesCode).
					Return(db.Permission{ID: test.ValidPermissionId, Code: ""}, nil)
				store.EXPECT().
					InsertUserPermission(gomock.Any(), arg2).
					Times(1).
					Return(db.UsersPermission{UserID: user.ID, PermissionID: test.ValidPermissionId}, nil)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)

				var gotUser db.User
				err := json.Unmarshal([]byte(rec.Body.String()), &gotUser)
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, user, gotUser, []string{"HashedPassword", "ID"}, "", "")
			},
		},
		{
			name: "Internal Server",
			body: map[string]any{
				"email":         user.Email,
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      password,
				"titleId":       user.TitleID,
			},
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				store.EXPECT().
					InsertUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
				require.Contains(t, rec.Body.String(), "sql: connection is already closed")
			},
		},
		{
			name: "Duplicate Username",
			body: map[string]any{
				"email":         user.Email,
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      password,
				"titleId":       user.TitleID,
			},
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				store.EXPECT().
					InsertUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrUniqueViolation)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
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

			srv, c, rec := setupHandlerTest(t, ctrl, tc, usersBasePath, string(data), http.MethodPut)

			err = srv.CreateUser(c)
			require.NoError(t, err)
			tc.checkRecorder(rec)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	updateUserParams := db.UpdateUserByIdParams{
		Email:          user1.Email,
		TitleID:        user1.TitleID,
		OgLanguageID:   user1.OgLanguageID,
		NewLanguageID:  user1.NewLanguageID,
		HashedPassword: user1.HashedPassword,
		ID:             user1.ID,
	}

	testCases := []testCase{
		{
			name:   "update user1 email",
			user:   user1,
			userId: user1.ID,
			body: `[
			{
				"op": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`,
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
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
			checkRecorder: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				var gotUser db.User
				err := json.Unmarshal([]byte(rec.Body.String()), &gotUser)
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, user1, gotUser, []string{"HashedPassword", "ID"}, "Email", "newemail2@email.com")
			},
		},
		{
			name:   "Id does not match",
			user:   user1,
			userId: user2.ID,
			body: `[
			{
				"op": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`,
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, rec.Code)
				require.Contains(t, rec.Body.String(), "Invalid user ID")
			},
		},
		{
			name:   "Invalid format for PatchUser",
			user:   user1,
			userId: user1.ID,
			body: `[
			{
				"wrong": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`,
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
				require.Contains(t, rec.Body.String(), "invalid operation {\"path\":\"/email\",\"value\":\"newemail2@email.com\",\"wrong\":\"replace\"}: unsupported operation")
			},
		},
		{
			name:   "User does not exist",
			user:   user2,
			userId: user2.ID,
			body: `[
			{
				"op": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`,
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				store.EXPECT().
					SelectUserById(gomock.Any(), user2.ID).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
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

			urlPath := usersBasePath + "/" + strconv.FormatInt(tc.user.ID, 10)
			srv, c, rec := setupHandlerTest(t, ctrl, tc, urlPath, tc.body.(string), http.MethodPatch)
			err := srv.UpdateUser(c, tc.userId)
			require.NoError(t, err)
			tc.checkRecorder(rec)
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
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				store.EXPECT().
					SelectUserByName(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					SelectUserPermissions(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(permissions, nil)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "UserNotFound",
			body: map[string]any{
				"username": "NotFound",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				store.EXPECT().
					SelectUserByName(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
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
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				store.EXPECT().
					SelectUserByName(gomock.Any(), gomock.Eq(user.Name)).
					Times(1).
					Return(user, nil)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
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
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
				store.EXPECT().
					SelectUserByName(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkRecorder: func(rec *httptest.ResponseRecorder) {
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

			srv, c, rec := setupHandlerTest(t, ctrl, tc, usersBasePath, string(data), http.MethodGet)

			err = srv.LoginUser(c)
			require.NoError(t, err)
			tc.checkRecorder(rec)
		})
	}
}

func TestCreateUserMiddleware(t *testing.T) {
	user, password := randomUser(t)
	testCases := []testCase{
		{
			name: "InvalidEmail",
			body: map[string]any{
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      password,
				"titleId":       user.TitleID,
				"email":         "invalid-email",
			},
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
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
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      "string",
				"titleId":       user.TitleID,
				"email":         user.Email,
			},
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX) {
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
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			urlPath := "/v1/users"
			ts, _ := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, data, ts, urlPath, http.MethodPost, "")

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}
