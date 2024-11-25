package api

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/test"
	"talkliketv.click/tltv/internal/util"
)

func TestFindTitles(t *testing.T) {
	t.Parallel()

	user, _ := randomUser(t)
	title := test.RandomTitle()
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
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
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
			buildStubs: func(stubs MockStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "{\"message\":\"parameter \\\"similarity\\\" in query has an error: value is required but missing\"}")
			},
			permissions: []string{db.ReadTitlesCode},
		},
		{
			name:   "Missing limit value",
			user:   user,
			values: map[string]any{"similarity": true, "limit": false},
			buildStubs: func(stubs MockStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "{\"message\":\"parameter \\\"limit\\\" in query has an error: value is required but missing\"}")
			},
			permissions: []string{db.ReadTitlesCode},
		},
	}

	for _, tc := range testCases {
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
			defer res.Body.Close()

			tc.checkResponse(res)
		})
	}
}

func TestAddTitle(t *testing.T) {
	user, _ := randomUser(t)
	title := test.RandomTitle()
	translate1 := randomTranslate(test.RandomPhrase(), title.OgLanguageID)
	translate2 := randomTranslate(test.RandomPhrase(), title.OgLanguageID)

	dbTranslates := []db.Translate{translate1, translate2}
	// remove directory after tests run
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			require.NoError(t, err)
		}
	}(test.AudioBasePath)
	err := os.MkdirAll(test.AudioBasePath, 0777)
	require.NoError(t, err)
	filename := test.AudioBasePath + "testAddTitle.txt"
	stringsSlice := []string{"This is the first sentence.", "This is the second sentence."}

	insertTitle := db.InsertTitleParams{
		Title:        title.Title,
		NumSubs:      2,
		OgLanguageID: title.OgLanguageID,
	}

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				data := []byte("This is the first sentence.\nThis is the second sentence.\n")
				fields := map[string]string{
					"titleName":  title.Title,
					"languageId": strconv.Itoa(int(title.OgLanguageID)),
				}
				return createMultiPartBody(t, data, filename, fields)
			},
			buildStubs: func(stubs MockStubs) {
				stubs.AudioFileX.EXPECT().
					GetLines(gomock.Any(), gomock.Any()).
					Return(stringsSlice, nil)
				stubs.MockQuerier.EXPECT().
					InsertTitle(gomock.Any(), insertTitle).
					Times(1).Return(title, nil)
				stubs.TranslateX.EXPECT().
					InsertNewPhrases(gomock.Any(), title, stubs.MockQuerier, stringsSlice).
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
				fields := map[string]string{
					"languageId": strconv.Itoa(int(title.OgLanguageID)),
				}
				return createMultiPartBody(t, data, filename, fields)
			},
			buildStubs: func(stubs MockStubs) {
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
				tooBigFile := test.AudioBasePath + "tooBigFile.txt"
				file, err := os.Create(tooBigFile)
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

				multiFile, err := os.Open(tooBigFile)
				require.NoError(t, err)
				body := new(bytes.Buffer)
				multiWriter := multipart.NewWriter(body)
				part, err := multiWriter.CreateFormFile("filePath", tooBigFile)
				require.NoError(t, err)
				_, err = io.Copy(part, multiFile)
				require.NoError(t, err)
				err = multiWriter.WriteField("titleName", title.Title)
				require.NoError(t, err)
				err = multiWriter.WriteField("languageId", strconv.Itoa(int(title.OgLanguageID)))
				require.NoError(t, err)
				require.NoError(t, multiWriter.Close())
				return body, multiWriter
			},
			buildStubs: func(stubs MockStubs) {
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
				fields := map[string]string{
					"titleName":  title.Title,
					"languageId": strconv.Itoa(int(title.OgLanguageID)),
				}
				return createMultiPartBody(t, data, filename, fields)
			},
			buildStubs: func(stubs MockStubs) {
				stubs.AudioFileX.EXPECT().
					GetLines(gomock.Any(), gomock.Any()).
					Return(stringsSlice, nil)
				stubs.MockQuerier.EXPECT().
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
			buildStubs: func(stubs MockStubs) {
			},
			multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				data := []byte("This is the first sentence.\nThis is the second sentence.\n")
				fields := map[string]string{
					"titleName":  title.Title,
					"languageId": strconv.Itoa(int(title.OgLanguageID)),
				}
				return createMultiPartBody(t, data, filename, fields)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				resBody := readBody(t, res)
				require.Contains(t, resBody, "security requirements failed: token claims don't match: provided claims do not match expected scopes")
			},
			permissions: []string{db.ReadTitlesCode},
		},
	}

	for _, tc := range testCases {
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
			defer res.Body.Close()

			tc.checkResponse(res)
		})
	}
}

func TestFindTitleById(t *testing.T) {
	t.Parallel()

	user, _ := randomUser(t)
	title := test.RandomTitle()

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
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
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlPath := titlesBasePath + "/" + strconv.FormatInt(title.ID, 10)
			//req, ts := setupServerTest(t, ctrl, tc, []byte(""), urlPath, http.MethodGet)
			ts, jwsToken := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, []byte(""), ts, urlPath, http.MethodGet, jwsToken)
			res, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			tc.checkResponse(res)
		})
	}
}

func TestDeleteTitleById(t *testing.T) {
	t.Parallel()

	user, _ := randomUser(t)
	title := test.RandomTitle()

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
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
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlPath := titlesBasePath + "/" + strconv.FormatInt(title.ID, 10)
			ts, jwsToken := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, []byte(""), ts, urlPath, http.MethodDelete, jwsToken)
			res, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			tc.checkResponse(res)
		})
	}
}

func TestTitlesTranslate(t *testing.T) {
	t.Parallel()

	user, _ := randomUser(t)
	title := test.RandomTitle()
	lang := randomLanguage()

	phrase1 := test.RandomPhrase()
	phrase2 := test.RandomPhrase()
	translate1 := db.Translate{
		PhraseID:   phrase1.Id,
		LanguageID: lang.ID,
		Phrase:     test.RandomString(8),
		PhraseHint: test.RandomString(8),
	}
	translate2 := db.Translate{
		PhraseID:   phrase2.Id,
		LanguageID: lang.ID,
		Phrase:     test.RandomString(8),
		PhraseHint: test.RandomString(8),
	}
	translates := []db.Translate{translate1, translate2}

	translatesReturn := []util.TranslatesReturn{
		{PhraseId: phrase1.Id,
			Text: test.RandomString(8)},
		{PhraseId: phrase2.Id, Text: test.RandomString(8)},
	}
	selectExistsTranslatesParams := db.SelectExistsTranslatesParams{
		LanguageID: lang.ID,
		ID:         title.ID,
	}
	selectTranslatesByTitleIdLangIdParams := db.SelectTranslatesByTitleIdLangIdParams{
		LanguageID: title.OgLanguageID,
		ID:         title.ID,
	}
	testCases := []testCase{
		{
			name: "OK",
			user: user,
			body: map[string]any{
				"newLanguageId": lang.ID,
				"titleId":       title.ID,
			},
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					SelectExistsTranslates(gomock.Any(), selectExistsTranslatesParams).
					Return(false, nil)
				stubs.MockQuerier.EXPECT().
					SelectTitleById(gomock.Any(), title.ID).
					Return(title, nil)
				stubs.MockQuerier.EXPECT().
					SelectLanguagesById(gomock.Any(), lang.ID).
					Return(lang, nil)
				stubs.MockQuerier.EXPECT().
					SelectTranslatesByTitleIdLangId(gomock.Any(), selectTranslatesByTitleIdLangIdParams).
					Return(translates, nil)
				stubs.TranslateX.EXPECT().
					TranslatePhrases(gomock.Any(), translates, lang).
					Return(translatesReturn, nil)
				stubs.TranslateX.EXPECT().
					InsertTranslates(gomock.Any(), stubs.MockQuerier, lang.ID, translatesReturn).
					Return(translates, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusCreated, res.StatusCode)
			},
			permissions: []string{"titles:w"},
		},
		{
			name: "Insert Translates Error",
			user: user,
			body: map[string]any{
				"newLanguageId": lang.ID,
				"titleId":       title.ID,
			},
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					SelectExistsTranslates(gomock.Any(), selectExistsTranslatesParams).
					Return(false, nil)
				stubs.MockQuerier.EXPECT().
					SelectTitleById(gomock.Any(), title.ID).
					Return(title, nil)
				stubs.MockQuerier.EXPECT().
					SelectLanguagesById(gomock.Any(), lang.ID).
					Return(lang, nil)
				stubs.MockQuerier.EXPECT().
					SelectTranslatesByTitleIdLangId(gomock.Any(), selectTranslatesByTitleIdLangIdParams).
					Return(translates, nil)
				stubs.TranslateX.EXPECT().
					TranslatePhrases(gomock.Any(), translates, lang).
					Return(translatesReturn, nil)
				stubs.TranslateX.EXPECT().
					InsertTranslates(gomock.Any(), stubs.MockQuerier, lang.ID, translatesReturn).
					Return([]db.Translate{}, sql.ErrConnDone)
				stubs.MockQuerier.EXPECT().
					DeleteTranslatesByLanguageId(
						gomock.Any(),
						db.DeleteTranslatesByLanguageIdParams{
							LanguageID: lang.ID,
							TitleID:    title.ID,
						}).
					Return(nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusInternalServerError, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "sql: connection is already closed")
			},
			permissions: []string{"titles:w"},
		},
		{
			name: "empty translates return",
			user: user,
			body: map[string]any{
				"newLanguageId": lang.ID,
				"titleId":       title.ID,
			},
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					SelectExistsTranslates(gomock.Any(), selectExistsTranslatesParams).
					Return(false, nil)
				stubs.MockQuerier.EXPECT().
					SelectTitleById(gomock.Any(), title.ID).
					Return(title, nil)
				stubs.MockQuerier.EXPECT().
					SelectLanguagesById(gomock.Any(), lang.ID).
					Return(lang, nil)
				stubs.MockQuerier.EXPECT().
					SelectTranslatesByTitleIdLangId(gomock.Any(), selectTranslatesByTitleIdLangIdParams).
					Return(translates, nil)
				stubs.TranslateX.EXPECT().
					TranslatePhrases(gomock.Any(), translates, lang).
					Return([]util.TranslatesReturn{}, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusInternalServerError, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "something went wrong")
			},
			permissions: []string{"titles:w"},
		},
		{
			name: "translates return error",
			user: user,
			body: map[string]any{
				"newLanguageId": lang.ID,
				"titleId":       title.ID,
			},
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					SelectExistsTranslates(gomock.Any(), selectExistsTranslatesParams).
					Return(false, nil)
				stubs.MockQuerier.EXPECT().
					SelectTitleById(gomock.Any(), title.ID).
					Return(title, nil)
				stubs.MockQuerier.EXPECT().
					SelectLanguagesById(gomock.Any(), lang.ID).
					Return(lang, nil)
				stubs.MockQuerier.EXPECT().
					SelectTranslatesByTitleIdLangId(gomock.Any(), selectTranslatesByTitleIdLangIdParams).
					Return(translates, nil)
				stubs.TranslateX.EXPECT().
					TranslatePhrases(gomock.Any(), translates, lang).
					Return([]util.TranslatesReturn{}, errors.New("new error for testing"))
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusInternalServerError, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "new error for testing")
			},
			permissions: []string{"titles:w"},
		},
		{
			name: "Already Exists",
			user: user,
			body: map[string]any{
				"newLanguageId": lang.ID,
				"titleId":       title.ID,
			},
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					SelectExistsTranslates(gomock.Any(), selectExistsTranslatesParams).
					Return(true, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "title already exists in that language")
			},
			permissions: []string{"titles:w"},
		},
		{
			name: "title not found",
			user: user,
			body: map[string]any{
				"newLanguageId": lang.ID,
				"titleId":       title.ID,
			},
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					SelectExistsTranslates(gomock.Any(), selectExistsTranslatesParams).
					Return(false, nil)
				stubs.MockQuerier.EXPECT().
					SelectTitleById(gomock.Any(), title.ID).
					Return(db.Title{}, sql.ErrNoRows)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "invalid title id")
			},
			permissions: []string{"titles:w"},
		},
		{
			name: "wrong permission",
			user: user,
			body: map[string]any{
				"newLanguageId": lang.ID,
				"titleId":       title.ID,
			},
			buildStubs: func(stubs MockStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "security requirements failed: token claims don't match: provided claims do not match expected scopes")
			},
			permissions: []string{"titles:r"},
		},
		{
			name: "bad request body",
			user: user,
			body: map[string]any{
				"newLanguage": lang.ID,
				"titleId":     title.ID,
			},
			buildStubs: func(stubs MockStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "request body has an error: doesn't match schema #/components/schemas/TitlesTranslateRequest: Error at \\\"/newLanguageId\\\"")
			},
			permissions: []string{"titles:w"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlPath := titlesBasePath + "/translate"
			ts, jwsToken := setupServerTest(t, ctrl, tc)
			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req := jsonRequest(t, data, ts, urlPath, http.MethodPost, jwsToken)
			res, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			tc.checkResponse(res)
		})
	}
}
