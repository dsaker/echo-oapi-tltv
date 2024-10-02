package api

import (
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"strconv"
	mockdb "talkliketv.click/tltv/db/mock"
	db "talkliketv.click/tltv/db/sqlc"
	"testing"
)

func TestFindTitles(t *testing.T) {
	user, _ := randomUser(t)
	title := randomTitle()
	listTitleParams := db.ListTitlesParams{
		Similarity: "similar",
		Limit:      10,
	}

	listTitleRow := db.ListTitlesRow{
		ID:           title.ID,
		Title:        title.Title,
		Similarity:   0.13333334,
		NumSubs:      title.NumSubs,
		OgLanguageID: title.OgLanguageID,
	}

	listTitlesRow := []db.ListTitlesRow{listTitleRow}

	testCases := []testCase{
		{
			name:   "OK",
			user:   user,
			values: map[string]any{"similarity": true, "limit": true},
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					ListTitles(gomock.Any(), listTitleParams).
					Times(1).
					Return(listTitlesRow, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var gotTitlesRow []db.ListTitlesRow
				err := json.Unmarshal([]byte(body), &gotTitlesRow)
				require.NoError(t, err)
				requireMatchAnyExcept(t, listTitlesRow[0], gotTitlesRow[0], nil, "", "")
			},
			permissions: []string{db.ReadTitlesCode},
		},
		{
			name:   "Missing similarity value",
			user:   user,
			values: map[string]any{"similarity": false, "limit": true},
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "{\"message\":\"parameter \\\"similarity\\\" in query has an error: value is required but missing\"}")
			},
			permissions: []string{db.ReadTitlesCode},
		},
		{
			name:   "Missing permission",
			user:   user,
			values: map[string]any{"similarity": false, "limit": true},
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "\"message\":\"security requirements failed: token claims don't match: provided claims do not match expected scopes\"")

			},
			permissions: []string{},
		},
		{
			name:   "Missing limit value",
			user:   user,
			values: map[string]any{"similarity": true, "limit": false},
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "{\"message\":\"parameter \\\"limit\\\" in query has an error: value is required but missing\"}")
			},
			permissions: []string{db.ReadTitlesCode},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, ts := setupServerTest(t, ctrl, tc, []byte(""), titlesBasePath, http.MethodGet)

			q := req.URL.Query()
			if tc.values["similarity"] == true {
				q.Add("similarity", "similar")
			}
			if tc.values["limit"] == true {
				q.Add("limit", "10")
			}
			req.URL.RawQuery = q.Encode()

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}

func TestAddTitle(t *testing.T) {
	user, _ := randomUser(t)

	title := randomTitle()

	insertTitle := db.InsertTitleParams{
		Title:        title.Title,
		NumSubs:      title.NumSubs,
		OgLanguageID: title.OgLanguageID,
	}

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			body: map[string]any{
				"numSubs":      title.NumSubs,
				"ogLanguageId": title.OgLanguageID,
				"title":        title.Title,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					InsertTitle(gomock.Any(), insertTitle).
					Times(1).
					Return(title, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var gotTitle db.Title
				err := json.Unmarshal([]byte(body), &gotTitle)
				require.NoError(t, err)
				requireMatchAnyExcept(t, title, gotTitle, nil, "", "")
			},
			permissions: []string{db.WriteTitlesCode},
		},
		{
			name: "Bad Request Body",
			user: user,
			body: map[string]any{
				"numSubs":    title.NumSubs,
				"ogLanguage": title.OgLanguageID,
				"title":      title.Title,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "request body has an error: doesn't match schema #/components/schemas/NewTitle: Error at ")
			},
			permissions: []string{db.WriteTitlesCode},
		},
		{
			name: "db connection closed",
			user: user,
			body: map[string]any{
				"numSubs":      title.NumSubs,
				"ogLanguageId": title.OgLanguageID,
				"title":        title.Title,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					InsertTitle(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Title{}, sql.ErrConnDone)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusInternalServerError, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "sql: connection is already closed")
			},
			permissions: []string{db.WriteTitlesCode},
		},
		{
			name: "missing permission",
			user: user,
			body: map[string]any{
				"numSubs":      title.NumSubs,
				"ogLanguageId": title.OgLanguageID,
				"title":        title.Title,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "\"message\":\"security requirements failed: token claims don't match: provided claims do not match expected scopes\"")
			},
			permissions: []string{},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, ts := setupServerTest(t, ctrl, tc, data, titlesBasePath, http.MethodPost)

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}

func TestFindTitleById(t *testing.T) {
	user, _ := randomUser(t)
	title := randomTitle()

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					SelectTitleById(gomock.Any(), title.ID).
					Times(1).
					Return(title, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var gotTitle db.Title
				err := json.Unmarshal([]byte(body), &gotTitle)
				require.NoError(t, err)
				requireMatchAnyExcept(t, title, gotTitle, nil, "", "")
			},
			permissions: []string{},
		},
		{
			name: "id not found",
			user: user,
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					SelectTitleById(gomock.Any(), title.ID).
					Times(1).
					Return(db.Title{}, sql.ErrNoRows)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "sql: no rows in result set")
			},
			permissions: []string{},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlPath := titlesBasePath + "/" + strconv.FormatInt(title.ID, 10)
			req, ts := setupServerTest(t, ctrl, tc, []byte(""), urlPath, http.MethodGet)

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}

func TestDeleteTitleById(t *testing.T) {
	user, _ := randomUser(t)
	title := randomTitle()

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					DeleteTitleById(gomock.Any(), title.ID).
					Times(1).Return(nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusNoContent, res.StatusCode)
			},
			permissions: []string{},
		},
		{
			name: "id not found",
			user: user,
			buildStubs: func(store *mockdb.MockQuerier) {
				store.EXPECT().
					DeleteTitleById(gomock.Any(), title.ID).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "sql: no rows in result set")
			},
			permissions: []string{},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlPath := titlesBasePath + "/" + strconv.FormatInt(title.ID, 10)
			req, ts := setupServerTest(t, ctrl, tc, []byte(""), urlPath, http.MethodDelete)

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}
