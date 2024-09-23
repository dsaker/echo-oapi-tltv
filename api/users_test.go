package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"reflect"
	"strconv"
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

func TestGetUser(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	testCases := []struct {
		name          string
		user          db.User
		userId        int64
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(res *http.Response)
	}{
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
			checkResponse: func(res *http.Response) {
				body := readBody(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
				var gotUser db.User
				err := json.Unmarshal([]byte(body), &gotUser)

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
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "Invalid user ID")
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
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "Error selecting user by id: sql: no rows in result set")
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

			ts, fa := newTestServer(t, store)

			jwsToken, err := fa.CreateJWSWithClaims([]string{db.ReadTitlesCode}, tc.user)
			require.NoError(t, err)
			urlPath := "/users/" + strconv.FormatInt(tc.userId, 10)
			res := request(t, []byte{}, ts, urlPath, "GET", string(jwsToken))
			require.NoError(t, err)
			tc.checkResponse(res)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	testCases := []struct {
		name          string
		user          db.User
		userId        int64
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(res *http.Response)
	}{
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
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusNoContent, res.StatusCode)
			},
		},
		{
			name:   "Id's don't match",
			user:   user1,
			userId: user2.ID,
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "Invalid user ID")
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
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "Error deleting user: sql: no rows in result set")
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

			ts, fa := newTestServer(t, store)

			jwsToken, err := fa.CreateJWSWithClaims([]string{db.ReadTitlesCode}, tc.user)
			require.NoError(t, err)
			urlPath := "/users/" + strconv.FormatInt(tc.userId, 10)
			res := request(t, []byte{}, ts, urlPath, "DELETE", string(jwsToken))
			require.NoError(t, err)
			tc.checkResponse(res)
		})
	}
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		body          map[string]any
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(res *http.Response)
	}{
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
			checkResponse: func(res *http.Response) {
				body := readBody(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
				var gotUser db.User
				err := json.Unmarshal([]byte(body), &gotUser)

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
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusInternalServerError, res.StatusCode)
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
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "duplicate key violation")
			},
		},
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
				body := readBody(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
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

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			ts, _ := newTestServer(t, store)
			urlPath := "/users"
			res, err := ts.Client().Post(ts.URL+urlPath, "application/json", bytes.NewReader(data))
			require.NoError(t, err)
			tc.checkResponse(res)
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

	testCases := []struct {
		name          string
		user          db.User
		userId        int64
		body          []byte
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(res *http.Response)
	}{
		{
			name:   "update user1 email",
			user:   user1,
			userId: user1.ID,
			body: []byte(`[
			{
				"op": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`),
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
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var gotUser db.User
				err := json.Unmarshal([]byte(body), &gotUser)
				require.NoError(t, err)
				requireMatchAnyExcept(t, user1, gotUser, []string{"HashedPassword", "ID"}, "Email", "newemail2@email.com")
			},
		},
		{
			name:   "Id does not match",
			user:   user1,
			userId: user2.ID,
			body: []byte(`[
			{
				"op": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`),
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "Invalid user ID")
			},
		},
		{
			name:   "Invalid format for PatchUser",
			user:   user1,
			userId: user1.ID,
			body: []byte(`[
			{
				"wrong": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`),
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "Error at \"/0\": property \"wrong\" is unsupported\n")
			},
		},
		{
			name:   "User does not exist",
			user:   user2,
			userId: user2.ID,
			body: []byte(`[
			{
				"op": "replace",
				"path": "/email",
				"value": "newemail2@email.com"
			}
		]`),
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					SelectUserById(gomock.Any(), user2.ID).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(res *http.Response) {
				body := readBody(t, res)
				require.Contains(t, body, "Error selecting user by id: pq:")
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

			ts, fa := newTestServer(t, store)

			jwsToken, err := fa.CreateJWSWithClaims([]string{db.ReadTitlesCode}, tc.user)
			if err != nil {
				log.Fatalln("error creating reader JWS:", err)
			}
			urlPath := "/users/" + strconv.FormatInt(tc.userId, 10)
			res := request(t, tc.body, ts, urlPath, "PATCH", string(jwsToken))
			require.NoError(t, err)
			tc.checkResponse(res)
		})
	}
}

func TestLoginUser(t *testing.T) {
	user, password := randomUser(t)
	permissions := []string{db.ReadTitlesCode}

	testCases := []struct {
		name          string
		body          map[string]any
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(res *http.Response)
	}{
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
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
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
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusNotFound, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "Invalid username or password")
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
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "Invalid username or password")
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
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusInternalServerError, res.StatusCode)
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

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			ts, _ := newTestServer(t, store)
			urlPath := "/users/login"
			res, err := ts.Client().Post(ts.URL+urlPath, "application/json", bytes.NewReader(data))
			require.NoError(t, err)
			tc.checkResponse(res)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:             util.RandomInt64(1, 1000),
		Name:           util.RandomName(),
		Email:          util.RandomEmail(),
		TitleID:        util.ValidTitleId,
		Flipped:        false,
		OgLanguageID:   util.ValidOgLanguageId,
		NewLanguageID:  util.ValidNewLanguageId,
		HashedPassword: hashedPassword,
	}
	return
}
