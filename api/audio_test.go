package api

import (
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
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/test"
	"testing"
)

func TestAudioFromFile(t *testing.T) {

	user, _ := randomUser(t)
	title := RandomTitle()
	translate1 := randomTranslate(test.RandomPhrase(), title.OgLanguageID)
	translate2 := randomTranslate(test.RandomPhrase(), title.OgLanguageID)

	fromLang := randomLanguage()
	toLang := randomLanguage()
	dbTranslates := []db.Translate{translate1, translate2}

	filename := "/tmp/sentences2.txt"
	stringsSlice := []string{"This is the first sentence.", "This is the second sentence."}

	audioFromTitleRequest := oapi.AudioFromTitleJSONRequestBody{
		FromLanguageId: fromLang.ID,
		TitleId:        title.ID,
		ToLanguageId:   toLang.ID,
	}

	//create a base path for storing mp3 audio files
	tmpAudioBasePath := "/tmp/audio/" +
		strconv.Itoa(int(title.ID)) + "/" +
		strconv.Itoa(int(title.OgLanguageID)) + "/"

	fromAudioBasePath := fmt.Sprintf("%s%d/", tmpAudioBasePath, fromLang.ID)
	toAudioBasePath := fmt.Sprintf("%s%d/", tmpAudioBasePath, toLang.ID)

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
				err = writer.WriteField("fileLanguageId", strconv.Itoa(int(title.OgLanguageID)))
				err = writer.WriteField("fromLanguageId", strconv.Itoa(int(lang1.ID)))
				err = writer.WriteField("toLanguageId", strconv.Itoa(int(lang2.ID)))
				err = writer.WriteField("toLanguageId", title.Title)
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
				// SelectLanguagesById(ctx context.Context, id int16) (Language, error)
				stubs.store.EXPECT().
					SelectLanguagesById(gomock.Any(), fromLang.ID).
					Return(fromLang, nil)
				stubs.store.EXPECT().
					SelectLanguagesById(gomock.Any(), toLang.ID).
					Return(toLang, nil)
				stubs.tr.EXPECT().
					CreateGoogleTranslateClient(gomock.Any()).
					Return(stubs.trc)
				stubs.tr.EXPECT().
					CreateGoogleTTSClient(gomock.Any()).
					Return(stubs.ttsc)
				// CreateTTS(echo.Context, db.Querier, clients.TTSClientX, clients.TranslateClientX, db.Language, db.Title, string) error
				stubs.tr.EXPECT().
					CreateTTS(gomock.Any(), stubs.store, stubs.ttsc, stubs.trc, fromLang, title, fromAudioBasePath).
					Return(nil)
				stubs.tr.EXPECT().
					CreateGoogleTranslateClient(gomock.Any()).
					Return(stubs.trc)
				stubs.tr.EXPECT().
					CreateGoogleTTSClient(gomock.Any()).
					Return(stubs.ttsc)
				stubs.tr.EXPECT().
					CreateTTS(gomock.Any(), stubs.store, stubs.ttsc, stubs.trc, toLang, title, toAudioBasePath).
					Return(nil)

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
		//{
		//	name: "Bad Request Body",
		//	user: user,
		//	multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
		//		data := []byte("This is the first sentence.\nThis is the second sentence.\n")
		//		err := os.WriteFile(filename, data, 0644)
		//		file, err := os.Open(filename)
		//		body := new(bytes.Buffer)
		//		writer := multipart.NewWriter(body)
		//		part, err := writer.CreateFormFile("filePath", filename)
		//		require.NoError(t, err)
		//		_, err = io.Copy(part, file)
		//		require.NoError(t, err)
		//		err = writer.WriteField("languageId", strconv.Itoa(int(title.OgLanguageID)))
		//		require.NoError(t, writer.Close())
		//		return body, writer
		//	},
		//	buildStubs: func(stubs buildStubs) {
		//	},
		//	checkResponse: func(res *http.Response) {
		//		require.Equal(t, http.StatusBadRequest, res.StatusCode)
		//		resBody := readBody(t, res)
		//		require.Contains(t, resBody, "{\"message\":\"request body has an error: doesn't match schema: Error at \\\"/titleName\\\": property")
		//	},
		//	permissions: []string{db.WriteTitlesCode},
		//},
		//{
		//	name: "File Too Big",
		//	user: user,
		//	multipartBody: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
		//		file, err := os.Create(filename)
		//		require.NoError(t, err)
		//		defer file.Close()
		//		writer := bufio.NewWriter(file)
		//		for i := 0; i < 64100; i++ {
		//			// Write random characters to the file
		//			char := byte('a')
		//			err = writer.WriteByte(char)
		//			require.NoError(t, err)
		//		}
		//		writer.Flush()
		//
		//		multiFile, err := os.Open(filename)
		//		body := new(bytes.Buffer)
		//		multiWriter := multipart.NewWriter(body)
		//		part, err := multiWriter.CreateFormFile("filePath", filename)
		//		require.NoError(t, err)
		//		_, err = io.Copy(part, multiFile)
		//		require.NoError(t, err)
		//		err = multiWriter.WriteField("titleName", title.Title)
		//		err = multiWriter.WriteField("languageId", strconv.Itoa(int(title.OgLanguageID)))
		//		require.NoError(t, multiWriter.Close())
		//		return body, multiWriter
		//	},
		//	buildStubs: func(stubs buildStubs) {
		//	},
		//	checkResponse: func(res *http.Response) {
		//		require.Equal(t, http.StatusBadRequest, res.StatusCode)
		//		resBody := readBody(t, res)
		//		require.Contains(t, resBody, "file too large")
		//	},
		//	permissions: []string{db.WriteTitlesCode},
		//},
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
			req, err := http.NewRequest(http.MethodPost, ts.URL+audioBasePath+"fromfile", multiBody)
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
