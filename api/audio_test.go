package api

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/test"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

func TestAudioFromTitle(t *testing.T) {

	user, _ := randomUser(t)
	title := test.RandomTitle()
	translate1 := randomTranslate(util.RandomPhrase(), title.OgLanguageID)
	translate2 := randomTranslate(util.RandomPhrase(), title.OgLanguageID)

	phraseIDs := []int64{translate1.PhraseID, translate2.PhraseID}
	fromLang := randomLanguage()
	toLang := randomLanguage()
	//dbTranslates := []db.Translate{translate1, translate2}

	//create a base path for storing mp3 audio files
	// TODO delete in cleanup
	tmpAudioBasePath := test.AudioBasePath + strconv.Itoa(int(title.ID)) + "/"
	err := os.MkdirAll(tmpAudioBasePath, 0777)
	require.NoError(t, err)

	filename := tmpAudioBasePath + "TestAudioFromTitle.txt"

	silenceBasePath := test.AudioBasePath + "silence/4SecSilence.mp3"
	fromAudioBasePath := fmt.Sprintf("%s%d/", tmpAudioBasePath, fromLang.ID)
	toAudioBasePath := fmt.Sprintf("%s%d/", tmpAudioBasePath, toLang.ID)

	testCases := []testCase{
		{
			name: "OK",
			body: map[string]any{
				"fromLanguageId": fromLang.ID,
				"titleId":        title.ID,
				"toLanguageId":   toLang.ID,
			},
			user: user,
			buildStubs: func(stubs MockStubs) {
				file, err := os.Create(filename)
				require.NoError(t, err)
				defer file.Close()
				stubs.MockQuerier.EXPECT().
					SelectTitleById(gomock.Any(), title.ID).
					Return(title, nil)
				// SelectLanguagesById(ctx context.Context, id int16) (Language, error)
				stubs.MockQuerier.EXPECT().
					SelectLanguagesById(gomock.Any(), fromLang.ID).
					Return(fromLang, nil)
				stubs.MockQuerier.EXPECT().
					SelectLanguagesById(gomock.Any(), toLang.ID).
					Return(toLang, nil)
				// CreateTTSForLang(echo.Context, db.Querier, db.Language, db.Title, string) error
				stubs.TranslateX.EXPECT().
					CreateTTS(gomock.Any(), stubs.MockQuerier, fromLang, title, fromAudioBasePath).
					Return(nil)
				stubs.TranslateX.EXPECT().
					CreateTTS(gomock.Any(), stubs.MockQuerier, toLang, title, toAudioBasePath).
					Return(nil)
				// SelectPhraseIdsByTitleId(ctx context.Context, titleID int64) ([]int64, error)
				stubs.MockQuerier.EXPECT().
					SelectPhraseIdsByTitleId(gomock.Any(), title.ID).
					Return(phraseIDs, nil)
				// BuildAudioInputFiles(echo.Context, []int64, db.Title, string, string, string, string) error
				stubs.AudioFileX.EXPECT().
					BuildAudioInputFiles(gomock.Any(), phraseIDs, title, silenceBasePath, fromAudioBasePath, toAudioBasePath, gomock.Any()).
					Return(nil)
				// CreateMp3Zip(echo.Context, db.Title, string) (*os.File, error)
				stubs.AudioFileX.EXPECT().
					CreateMp3ZipWithFfmpeg(gomock.Any(), title, gomock.Any()).
					Return(file, nil)

			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
			},
			permissions: []string{db.WriteTitlesCode},
		},
		{
			name: "Bad Request Body",
			body: map[string]any{
				"fromLanguageId": fromLang.ID,
				"titleId":        user.ID,
				"toLanguage":     toLang.ID,
			},
			user: user,
			buildStubs: func(stubs MockStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				resBody := readBody(t, res)
				require.Contains(t, resBody, "{\"message\":\"request body has an error: doesn't match schema: Error at \\\"/toLanguageId\\\": property \\\"toLanguageId\\\" is missing\"}")
			},
			permissions: []string{db.WriteTitlesCode},
		},
		{
			name: "Db connection closed",
			body: map[string]any{
				"fromLanguageId": fromLang.ID,
				"titleId":        title.ID,
				"toLanguageId":   toLang.ID,
			},
			user: user,
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					SelectTitleById(gomock.Any(), title.ID).
					Return(db.Title{}, sql.ErrConnDone)
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
			body: map[string]any{
				"fromLanguageId": fromLang.ID,
				"titleId":        title.ID,
				"toLanguage":     toLang.ID,
			},
			buildStubs: func(stubs MockStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				resBody := readBody(t, res)
				require.Contains(t, resBody, "security requirements failed: token claims don't match: provided claims do not match expected scopes")
			},
			permissions: []string{db.ReadTitlesCode},
			cleanUp: func(t *testing.T) {
				err = os.Remove(tmpAudioBasePath)
				require.NoError(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts, jwsToken := setupServerTest(t, ctrl, tc)
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req := jsonRequest(t, data, ts, audioBasePath+"/fromtitle", http.MethodPost, jwsToken)
			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
			require.NoError(t, err)
		})
	}
}

func TestAudioFromFile(t *testing.T) {

	user, _ := randomUser(t)
	title := test.RandomTitle()
	translate1 := randomTranslate(util.RandomPhrase(), title.OgLanguageID)
	translate2 := randomTranslate(util.RandomPhrase(), title.OgLanguageID)

	phraseIDs := []int64{translate1.PhraseID, translate2.PhraseID}
	fromLang := randomLanguage()
	toLang := randomLanguage()
	dbTranslates := []db.Translate{translate1, translate2}

	//create a base path for storing mp3 audio files
	// TODO delete in cleanup
	tmpAudioBasePath := test.AudioBasePath + strconv.Itoa(int(title.ID)) + "/"
	err := os.MkdirAll(tmpAudioBasePath, 0777)
	require.NoError(t, err)

	filename := tmpAudioBasePath + "TestAudioFromFile.txt"
	stringsSlice := []string{"This is the first sentence.", "This is the second sentence."}

	silenceBasePath := test.AudioBasePath + "silence/4SecSilence.mp3"
	fromAudioBasePath := fmt.Sprintf("%s%d/", tmpAudioBasePath, fromLang.ID)
	toAudioBasePath := fmt.Sprintf("%s%d/", tmpAudioBasePath, toLang.ID)

	insertTitle := db.InsertTitleParams{
		Title:        title.Title,
		NumSubs:      2,
		OgLanguageID: title.OgLanguageID,
	}

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			buildStubs: func(stubs MockStubs) {
				file, err := os.Create(filename)
				require.NoError(t, err)
				defer file.Close()
				// GetLines(echo.Context, multipart.File) ([]string, error)
				stubs.AudioFileX.EXPECT().
					GetLines(gomock.Any(), gomock.Any()).
					Return(stringsSlice, nil)
				stubs.MockQuerier.EXPECT().
					InsertTitle(gomock.Any(), insertTitle).
					Times(1).Return(title, nil)
				stubs.TranslateX.EXPECT().
					InsertNewPhrases(gomock.Any(), title, stubs.MockQuerier, stringsSlice).
					Times(1).Return(dbTranslates, nil)
				// SelectLanguagesById(ctx context.Context, id int16) (Language, error)
				stubs.MockQuerier.EXPECT().
					SelectLanguagesById(gomock.Any(), fromLang.ID).
					Return(fromLang, nil)
				stubs.MockQuerier.EXPECT().
					SelectLanguagesById(gomock.Any(), toLang.ID).
					Return(toLang, nil)
				// CreateTTSForLang(echo.Context, db.Querier, db.Language, db.Title, string) error
				stubs.TranslateX.EXPECT().
					CreateTTS(gomock.Any(), stubs.MockQuerier, fromLang, title, fromAudioBasePath).
					Return(nil)
				stubs.TranslateX.EXPECT().
					CreateTTS(gomock.Any(), stubs.MockQuerier, toLang, title, toAudioBasePath).
					Return(nil)
				// SelectPhraseIdsByTitleId(ctx context.Context, titleID int64) ([]int64, error)
				stubs.MockQuerier.EXPECT().
					SelectPhraseIdsByTitleId(gomock.Any(), title.ID).
					Return(phraseIDs, nil)
				// BuildAudioInputFiles(echo.Context, []int64, db.Title, string, string, string, string) error
				stubs.AudioFileX.EXPECT().
					BuildAudioInputFiles(gomock.Any(), phraseIDs, title, silenceBasePath, fromAudioBasePath, toAudioBasePath, gomock.Any()).
					Return(nil)
				// CreateMp3Zip(echo.Context, db.Title, string) (*os.File, error)
				stubs.AudioFileX.EXPECT().
					CreateMp3ZipWithFfmpeg(gomock.Any(), title, gomock.Any()).
					Return(file, nil)

			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
			},
			multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				data := []byte("This is the first sentence.\nThis is the second sentence.\n")

				formMap := map[string]string{
					"fileLanguageId": strconv.Itoa(int(title.OgLanguageID)),
					"fromLanguageId": strconv.Itoa(int(fromLang.ID)),
					"toLanguageId":   strconv.Itoa(int(toLang.ID)),
					"titleName":      title.Title}
				return createMultiPartBody(t, data, filename, formMap)
			},
			permissions: []string{db.WriteTitlesCode},
		},
		{
			name: "Bad Request Body",
			user: user,
			multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				data := []byte("This is the first sentence.\nThis is the second sentence.\n")
				formMap := map[string]string{
					"fileLanguageId": strconv.Itoa(int(title.OgLanguageID)),
					"fromLanguageId": strconv.Itoa(int(fromLang.ID)),
					"toLanguageId":   strconv.Itoa(int(toLang.ID)),
				}
				return createMultiPartBody(t, data, filename, formMap)
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
				body := new(bytes.Buffer)
				multiWriter := multipart.NewWriter(body)
				part, err := multiWriter.CreateFormFile("filePath", tooBigFile)
				require.NoError(t, err)
				_, err = io.Copy(part, multiFile)
				require.NoError(t, err)
				fieldMap := map[string]string{
					"fileLanguageId": strconv.Itoa(int(title.OgLanguageID)),
					"fromLanguageId": strconv.Itoa(int(fromLang.ID)),
					"toLanguageId":   strconv.Itoa(int(toLang.ID)),
					"titleName":      title.Title}
				for field, value := range fieldMap {
					err = multiWriter.WriteField(field, value)
					require.NoError(t, err)
				}
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

				formMap := map[string]string{
					"fileLanguageId": strconv.Itoa(int(title.OgLanguageID)),
					"fromLanguageId": strconv.Itoa(int(fromLang.ID)),
					"toLanguageId":   strconv.Itoa(int(toLang.ID)),
					"titleName":      title.Title}
				return createMultiPartBody(t, data, filename, formMap)
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
			multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				data := []byte("This is the first sentence.\nThis is the second sentence.\n")

				formMap := map[string]string{
					"fileLanguageId": strconv.Itoa(int(title.OgLanguageID)),
					"fromLanguageId": strconv.Itoa(int(fromLang.ID)),
					"toLanguageId":   strconv.Itoa(int(toLang.ID)),
					"titleName":      title.Title}
				return createMultiPartBody(t, data, filename, formMap)
			},
			buildStubs: func(stubs MockStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				resBody := readBody(t, res)
				require.Contains(t, resBody, "security requirements failed: token claims don't match: provided claims do not match expected scopes")
			},
			permissions: []string{db.ReadTitlesCode},
			cleanUp: func(t *testing.T) {
				err = os.Remove(tmpAudioBasePath)
				require.NoError(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts, jwsToken := setupServerTest(t, ctrl, tc)
			multiBody, multiWriter := tc.multipartBody(t)
			req, err := http.NewRequest(http.MethodPost, ts.URL+audioBasePath+"/fromfile", multiBody)
			require.NoError(t, err)

			req.Header.Set("Authorization", "Bearer "+string(jwsToken))

			req.Header.Set("Content-Type", multiWriter.FormDataContentType())
			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
			require.NoError(t, err)
		})
	}
}
