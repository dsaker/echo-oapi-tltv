package api

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	mockdb "talkliketv.click/tltv/db/mock"
	db "talkliketv.click/tltv/db/sqlc"
	"testing"
)

type titlesTestCase struct {
	name          string
	user          db.User
	body          db.ListTitlesParams
	similarity    bool
	limit         bool
	buildStubs    func(store *mockdb.MockQuerier)
	checkResponse func(res *http.Response)
}

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
		LanguageID:   title.LanguageID,
	}

	listTitlesRow := []db.ListTitlesRow{listTitleRow}

	testCases := []titlesTestCase{
		{
			name: "OK",
			user: user,
			body: db.ListTitlesParams{
				Similarity: "similar",
				Limit:      10,
			},
			similarity: true,
			limit:      true,
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
				requireMatchAnyExcept(t, listTitlesRow[0], gotTitlesRow[0], []string{}, "", "")
			},
		},
		{
			name: "Missing similarity value",
			user: user,
			body: db.ListTitlesParams{
				Similarity: "",
				Limit:      10,
			},
			similarity: false,
			limit:      true,
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "Invalid format for parameter similarity: query parameter 'similarity' is required")

			},
		},
		{
			name: "Missing limit value",
			user: user,
			body: db.ListTitlesParams{
				Similarity: "",
				Limit:      10,
			},
			similarity: true,
			limit:      false,
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "Invalid format for parameter limit: query parameter 'limit' is required")

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req := setupServerTest(t, ctrl, tc)
			//store := mockdb.NewMockQuerier(ctrl)
			//tc.buildStubs(store)
			//
			//spec, err := GetSwagger()
			//require.NoError(t, err)
			//
			//e := echo.New()
			//svr := NewServer(e, testCfg, store, spec)
			//
			//RegisterHandlersWithBaseURL(e, svr, "/v1")
			//ts := httptest.NewServer(e)
			//
			//jwsToken, err := svr.fa.CreateJWSWithClaims([]string{db.ReadTitlesCode}, user)
			//require.NoError(t, err)
			//urlPath := "/v1/titles"
			//
			//req := serverRequest(t, nil, ts, urlPath, http.MethodGet, string(jwsToken))
			q := req.URL.Query()
			if tc.similarity {
				q.Add("similarity", "similar")
			}
			if tc.limit {
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
		LanguageID:   title.LanguageID,
		OgLanguageID: title.OgLanguageID,
	}

	testCases := []struct {
		name          string
		body          interface{}
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(res *http.Response)
	}{
		{
			name: "OK",
			body: map[string]any{
				"languageId":   title.LanguageID,
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
			svr := NewServer(e, testCfg, store, spec)

			RegisterHandlersWithBaseURL(e, svr, "/v1")
			ts := httptest.NewServer(e)

			jwsToken, err := svr.fa.CreateJWSWithClaims([]string{db.ReadTitlesCode, db.WriteTitlesCode}, user)
			require.NoError(t, err)
			urlPath := "/v1/titles"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := serverRequest(t, data, ts, urlPath, http.MethodPost, string(jwsToken))

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}
