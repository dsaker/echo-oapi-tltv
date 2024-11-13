package api

import (
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/test"
	"testing"
)

func TestAddUserPermission(t *testing.T) {
	user, _ := randomUser(t)

	userPermission := db.UsersPermission{
		PermissionID: 1,
		UserID:       user.ID,
	}

	insertUsersPermission := db.InsertUserPermissionParams{
		UserID:       userPermission.UserID,
		PermissionID: userPermission.PermissionID,
	}

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			body: map[string]any{
				"permissionId": 1,
				"userId":       user.ID,
			},
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					InsertUserPermission(gomock.Any(), insertUsersPermission).
					Times(1).
					Return(userPermission, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusCreated, res.StatusCode)
				body := readBody(t, res)
				var got db.UsersPermission
				err := json.Unmarshal([]byte(body), &got)
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, userPermission, got, nil, "", "")
			},
			permissions: []string{db.GlobalAdminCode},
		},
		{
			name: "Bad Request Body",
			user: user,
			body: map[string]any{
				"permission": 1,
				"userId":     user.ID,
			},
			buildStubs: func(stubs MockStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "request body has an error: doesn't match schema #/components/schemas/NewUserPermission: Error at ")
			},
			permissions: []string{db.GlobalAdminCode},
		},
		{
			name: "MockQuerier connection closed",
			user: user,
			body: map[string]any{
				"permissionId": 1,
				"userId":       user.ID,
			},
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					InsertUserPermission(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.UsersPermission{}, sql.ErrConnDone)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusInternalServerError, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "sql: connection is already closed")
			},
			permissions: []string{db.GlobalAdminCode},
		},
		{
			name: "missing permission",
			user: user,
			body: map[string]any{
				"permissionId": 1,
				"userId":       user.ID,
			},
			buildStubs: func(stubs MockStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "\"message\":\"security requirements failed: token claims don't match: provided claims do not match expected scopes\"")
			},
			permissions: []string{},
		},
		{
			name: "foreign key violation",
			user: user,
			body: map[string]any{
				"permissionId": 1,
				"userId":       user.ID,
			},
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					InsertUserPermission(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.UsersPermission{}, db.ErrForeignKeyViolation)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "pq: insert or update on table \"users_permissions\" violates foreign key constraint \"users_permissions_user_id_fkey\"")
			},
			permissions: []string{db.GlobalAdminCode},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			ts, jwsToken := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, data, ts, usersPermissionBasePath, http.MethodPost, jwsToken)
			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}
