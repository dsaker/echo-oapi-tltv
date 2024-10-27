package api

import (
	"bufio"
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
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/test"
	"testing"
)

func TestFindTitles(t *testing.T) {
	user, _ := randomUser(t)
	title := RandomTitle()
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
			buildStubs: func(stubs buildStubs) {
				stubs.store.EXPECT().
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
				test.RequireMatchAnyExcept(t, listTitlesRow[0], gotTitlesRow[0], nil, "", "")
			},
			permissions: []string{db.ReadTitlesCode},
		},
		{
			name:   "Missing similarity value",
			user:   user,
			values: map[string]any{"similarity": false, "limit": true},
			buildStubs: func(stubs buildStubs) {
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
			buildStubs: func(stubs buildStubs) {
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
			buildStubs: func(stubs buildStubs) {
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

			//req, ts := setupServerTest(t, ctrl, tc, []byte(""), titlesBasePath, http.MethodGet)
			ts, jwsToken := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, []byte(""), ts, titlesBasePath, http.MethodGet, jwsToken)
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
	title := RandomTitle()
	translate1 := randomTranslate(test.RandomPhrase(), title.OgLanguageID)
	translate2 := randomTranslate(test.RandomPhrase(), title.OgLanguageID)

	dbTranslates := []db.Translate{translate1, translate2}

	filename := "/tmp/sentences2.txt"
	stringsSlice := []string{"This is the first sentence.", "This is the second sentence."}

	//create a base path for storing mp3 audio files
	audioBasePath := "/tmp/audio/" +
		strconv.Itoa(int(title.ID)) + "/" +
		strconv.Itoa(int(title.OgLanguageID)) + "/"

	insertTitle := db.InsertTitleParams{
		Title:        title.Title,
		NumSubs:      2,
		OgLanguageID: title.OgLanguageID,
	}

	filename = "/tmp/filename.txt"

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				data := []byte("This is the first sentence.\nThis is the second sentence.\n")
				err := os.WriteFile(filename, data, 0777)
				file, err := os.Open(filename)
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("filePath", filename)
				require.NoError(t, err)
				_, err = io.Copy(part, file)
				require.NoError(t, err)
				err = writer.WriteField("titleName", title.Title)
				err = writer.WriteField("languageId", strconv.Itoa(int(title.OgLanguageID)))
				require.NoError(t, writer.Close())
				return body, writer
			},
			buildStubs: func(stubs buildStubs) {
				stubs.store.EXPECT().
					InsertTitle(gomock.Any(), insertTitle).
					Times(1).Return(title, nil)
				stubs.tr.EXPECT().
					InsertNewPhrases(gomock.Any(), title, stubs.store, stringsSlice).
					Times(1).Return(dbTranslates, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var gotTitle db.Title
				err := json.Unmarshal([]byte(body), &gotTitle)
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, title, gotTitle, nil, "", "")
			},
			permissions: []string{db.WriteTitlesCode},
		},
		{
			name: "Bad Request Body",
			user: user,
			multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				data := []byte("This is the first sentence.\nThis is the second sentence.\n")
				err := os.WriteFile(filename, data, 0644)
				file, err := os.Open(filename)
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("filePath", filename)
				require.NoError(t, err)
				_, err = io.Copy(part, file)
				require.NoError(t, err)
				err = writer.WriteField("languageId", strconv.Itoa(int(title.OgLanguageID)))
				require.NoError(t, writer.Close())
				return body, writer
			},
			buildStubs: func(stubs buildStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				resBody := readBody(t, res)
				require.Contains(t, resBody, "{\"message\":\"request body has an error: doesn't match schema: Error at \\\"/titleName\\\": property")
			},
			permissions: []string{db.WriteTitlesCode},
		},
		{
			name: "File Too Big",
			user: user,
			multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				file, err := os.Create(filename)
				require.NoError(t, err)
				defer file.Close()
				writer := bufio.NewWriter(file)
				for i := 0; i < 64100; i++ {
					// Write random characters to the file
					char := byte('a')
					err = writer.WriteByte(char)
					require.NoError(t, err)
				}
				writer.Flush()

				multiFile, err := os.Open(filename)
				body := new(bytes.Buffer)
				multiWriter := multipart.NewWriter(body)
				part, err := multiWriter.CreateFormFile("filePath", filename)
				require.NoError(t, err)
				_, err = io.Copy(part, multiFile)
				require.NoError(t, err)
				err = multiWriter.WriteField("titleName", title.Title)
				err = multiWriter.WriteField("languageId", strconv.Itoa(int(title.OgLanguageID)))
				require.NoError(t, multiWriter.Close())
				return body, multiWriter
			},
			buildStubs: func(stubs buildStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				resBody := readBody(t, res)
				require.Contains(t, resBody, "file too large")
			},
			permissions: []string{db.WriteTitlesCode},
		},
		{
			name: "Db connection closed",
			user: user,
			multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				data := []byte("This is the first sentence.\nThis is the second sentence.\n")
				err := os.WriteFile(filename, data, 0644)
				file, err := os.Open(filename)
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("filePath", filename)
				require.NoError(t, err)
				_, err = io.Copy(part, file)
				require.NoError(t, err)
				err = writer.WriteField("titleName", title.Title)
				err = writer.WriteField("languageId", strconv.Itoa(int(title.OgLanguageID)))
				require.NoError(t, writer.Close())
				return body, writer
			},
			buildStubs: func(stubs buildStubs) {
				stubs.store.EXPECT().
					InsertTitle(gomock.Any(), insertTitle).
					Times(1).Return(db.Title{}, sql.ErrConnDone)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusInternalServerError, res.StatusCode)
				resBody := readBody(t, res)
				require.Contains(t, resBody, "sql: connection is already closed")
			},
			permissions: []string{db.WriteTitlesCode},
		},
		{
			name: "missing permission",
			user: user,
			multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				data := []byte("This is the first sentence.\nThis is the second sentence.\n")
				err := os.WriteFile(filename, data, 0644)
				file, err := os.Open(filename)
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("filePath", filename)
				require.NoError(t, err)
				_, err = io.Copy(part, file)
				require.NoError(t, err)
				err = writer.WriteField("titleName", title.Title)
				err = writer.WriteField("languageId", strconv.Itoa(int(title.OgLanguageID)))
				require.NoError(t, writer.Close())
				return body, writer
			},
			buildStubs: func(stubs buildStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				resBody := readBody(t, res)
				require.Contains(t, resBody, "security requirements failed: token claims don't match: provided claims do not match expected scopes")
			},
			permissions: []string{db.ReadTitlesCode},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts, jwsToken := setupServerTest(t, ctrl, tc)
			multiBody, multiWriter := tc.multipartBody(t)
			req, err := http.NewRequest(http.MethodPost, ts.URL+titlesBasePath, multiBody)
			require.NoError(t, err)

			req.Header.Set("Authorization", "Bearer "+string(jwsToken))

			req.Header.Set("Content-Type", multiWriter.FormDataContentType())
			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
			err = os.RemoveAll(audioBasePath)
			require.NoError(t, err)
			err = os.Remove(filename)
			require.NoError(t, err)
		})
	}
}

func TestFindTitleById(t *testing.T) {
	user, _ := randomUser(t)
	title := RandomTitle()

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			buildStubs: func(stubs buildStubs) {
				stubs.store.EXPECT().
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
				test.RequireMatchAnyExcept(t, title, gotTitle, nil, "", "")
			},
			permissions: []string{},
		},
		{
			name: "id not found",
			user: user,
			buildStubs: func(stubs buildStubs) {
				stubs.store.EXPECT().
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
			//req, ts := setupServerTest(t, ctrl, tc, []byte(""), urlPath, http.MethodGet)
			ts, jwsToken := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, []byte(""), ts, urlPath, http.MethodGet, jwsToken)
			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}

func TestDeleteTitleById(t *testing.T) {
	user, _ := randomUser(t)
	title := RandomTitle()

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			buildStubs: func(stubs buildStubs) {
				stubs.store.EXPECT().
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
			buildStubs: func(stubs buildStubs) {
				stubs.store.EXPECT().
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
			ts, jwsToken := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, []byte(""), ts, urlPath, http.MethodDelete, jwsToken)
			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}
