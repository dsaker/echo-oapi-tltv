package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	mockdb "talkliketv.click/tltv/db/mock"
	db "talkliketv.click/tltv/db/sqlc"
	mock "talkliketv.click/tltv/internal/mock"
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
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
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
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
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
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
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
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
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
	phrase1 := randomPhrase()
	phrase2 := randomPhrase()
	translate1 := randomTranslate(phrase1, validLanguageId)
	translate2 := randomTranslate(phrase2, validLanguageId)

	dbTranslates := []db.Translate{translate1, translate2}

	filename := "sentences.txt"
	stringsSlice := []string{"This is the first sentence.", "This is the second sentence."}
	//create base path for storing mp3 audio files
	audioBasePath := "/tmp/audio/" +
		strconv.Itoa(int(title.ID)) + "/" +
		strconv.Itoa(int(title.OgLanguageID)) + "/"

	insertTitle := db.InsertTitleParams{
		Title:        title.Title,
		NumSubs:      2,
		OgLanguageID: title.OgLanguageID,
	}

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			body: map[string]any{
				"titleName":  "sentences",
				"languageId": title.OgLanguageID,
				"filePath":   filename,
			},
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
				store.EXPECT().
					SelectLanguagesById(gomock.Any(), title.OgLanguageID).
					Times(1).Return(validLanguageModel, nil)
				store.EXPECT().
					InsertTitle(gomock.Any(), insertTitle).
					Times(1).Return(title, nil)
				text.EXPECT().
					InsertPhrases(gomock.Any(), title, store, stringsSlice, 2).
					Times(1).Return(dbTranslates, nil)
				text.EXPECT().
					TextToSpeech(gomock.Any(), dbTranslates, audioBasePath, validLanguageModel.Tag).
					Times(1).Return(nil)

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
		//{
		//	name: "Bad Request Body",
		//	user: user,
		//	body: map[string]any{
		//		"numSubs":    title.NumSubs,
		//		"ogLanguage": title.OgLanguageID,
		//		"title":      title.Title,
		//	},
		//	buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
		//	},
		//	checkResponse: func(res *http.Response) {
		//		require.Equal(t, http.StatusBadRequest, res.StatusCode)
		//		body := readBody(t, res)
		//		require.Contains(t, body, "request body has an error: doesn't match schema #/components/schemas/NewTitle: Error at ")
		//	},
		//	permissions: []string{db.WriteTitlesCode},
		//},
		//{
		//	name: "db connection closed",
		//	user: user,
		//	body: map[string]any{
		//		"numSubs":      title.NumSubs,
		//		"ogLanguageId": title.OgLanguageID,
		//		"title":        title.Title,
		//	},
		//	buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
		//		store.EXPECT().
		//			InsertTitle(gomock.Any(), gomock.Any()).
		//			Times(1).
		//			Return(db.Title{}, sql.ErrConnDone)
		//	},
		//	checkResponse: func(res *http.Response) {
		//		require.Equal(t, http.StatusInternalServerError, res.StatusCode)
		//		body := readBody(t, res)
		//		require.Contains(t, body, "sql: connection is already closed")
		//	},
		//	permissions: []string{db.WriteTitlesCode},
		//},
		//{
		//	name: "missing permission",
		//	user: user,
		//	body: map[string]any{
		//		"numSubs":      title.NumSubs,
		//		"ogLanguageId": title.OgLanguageID,
		//		"title":        title.Title,
		//	},
		//	buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
		//	},
		//	checkResponse: func(res *http.Response) {
		//		require.Equal(t, http.StatusForbidden, res.StatusCode)
		//		body := readBody(t, res)
		//		require.Contains(t, body, "\"message\":\"security requirements failed: token claims don't match: provided claims do not match expected scopes\"")
		//	},
		//	permissions: []string{},
		//},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			//data, err := json.Marshal(body)
			//require.NoError(t, err)

			//req, ts := setupServerTest(t, ctrl, tc, data, titlesBasePath, http.MethodPost)

			data := []byte("This is the first sentence.\nThis is the second sentence.\n")
			body := bytes.NewBuffer(data)
			writer := multipart.NewWriter(body)
			file, err := os.Open(filename)
			require.NoError(t, err)
			part, err := writer.CreateFormFile("filePath", filename)
			require.NoError(t, err)
			_, err = io.Copy(part, file)
			require.NoError(t, err)
			err = writer.WriteField("titleName", title.Title)
			err = writer.WriteField("languageId", strconv.Itoa(int(title.OgLanguageID)))
			require.NoError(t, writer.Close())

			//text := mock.NewMockTranslateX(ctrl)
			//store := mockdb.NewMockQuerier(ctrl)
			//tc.buildStubs(store, text)
			//
			//e, srv := NewServer(testCfg, store, text)
			//
			//ts := httptest.NewServer(e)
			//
			//jwsToken, err := srv.fa.CreateJWSWithClaims(tc.permissions, tc.user)
			//require.NoError(t, err)

			//req, err := http.NewRequest(http.MethodPost, ts.URL+titlesBasePath, body)
			//require.NoError(t, err)

			req, ts := setupServerTest(t, ctrl, tc, body, titlesBasePath, http.MethodPost)

			req.Header.Set("Authorization", "Bearer "+string(jwsToken))

			req.Header.Set("Content-Type", writer.FormDataContentType())
			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
			err = os.RemoveAll(audioBasePath)
			require.NoError(t, err)
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
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
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
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
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
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
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
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
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
