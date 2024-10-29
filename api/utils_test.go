package api

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	mocka "talkliketv.click/tltv/internal/mock/audiofile"
	mockc "talkliketv.click/tltv/internal/mock/clients"
	mockdb "talkliketv.click/tltv/internal/mock/db"
	mockt "talkliketv.click/tltv/internal/mock/translates"
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/test"
	"talkliketv.click/tltv/internal/token"
	"talkliketv.click/tltv/internal/util"
	"testing"
	"time"
)

var (
	testCfg config.Config
)

const (
	usersBasePath           = "/v1/users"
	audioBasePath           = "/v1/audio"
	titlesBasePath          = "/v1/titles"
	usersPermissionBasePath = "/v1/userspermissions"
	phrasesBasePath         = "/v1/phrases"
	usersPhrasesBasePath    = "/v1/usersphrases"
	languagesBasePath       = "/v1/languages"
	validLanguageId         = 27
	testAudioBasePath       = "../tmp/test/audio/"
)

type mockStubs struct {
	mockQuerier      *mockdb.MockQuerier
	translateX       *mockt.MockTranslateX
	translateClientX *mockc.MockTranslateClientX
	ttsClientX       *mockc.MockTTSClientX
	audioFileX       *mocka.MockAudioFileX
}

func newMockStubs(ctrl *gomock.Controller) mockStubs {
	return mockStubs{
		mockQuerier:      mockdb.NewMockQuerier(ctrl),
		translateX:       mockt.NewMockTranslateX(ctrl),
		translateClientX: mockc.NewMockTranslateClientX(ctrl),
		ttsClientX:       mockc.NewMockTTSClientX(ctrl),
		audioFileX:       mocka.NewMockAudioFileX(ctrl),
	}
}

type testCase struct {
	name          string
	body          interface{}
	user          db.User
	extraInt      int64
	buildStubs    func(stubs mockStubs)
	multipartBody func(t *testing.T) (*bytes.Buffer, *multipart.Writer)
	checkRecorder func(rec *httptest.ResponseRecorder)
	checkResponse func(res *http.Response)
	values        map[string]any
	permissions   []string
	cleanUp       func(*testing.T)
}

func TestMain(m *testing.M) {
	testCfg = SetTestConfigs()
	flag.Parse()
	os.Exit(m.Run())
}

func SetTestConfigs() (config config.Config) {

	// get port and debug from commandline flags... if not present use defaults
	flag.StringVar(&config.Port, "port", "8080", "API server port")

	flag.StringVar(&config.Env, "env", "development", "Environment (development|staging|cloud)")
	flag.DurationVar(&config.CtxTimeout, "ctx-timeout", 3*time.Second, "Context timeout for db queries in seconds")

	flag.StringVar(&config.TTSBasePath, "tts-base-path", testAudioBasePath, "text-to-speech base path temporary storage of mp3 audio files")

	flag.Int64Var(&config.FileUploadLimit, "upload-size-limit", 8, "File upload size limit in KB (default is 4)")
	return config

}

func readBody(t *testing.T, rs *http.Response) string {
	// Read the checkResponse body from the test server.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(rs.Body)

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return string(body)
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = test.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:             test.RandomInt64(),
		Name:           test.RandomString(8),
		Email:          test.RandomEmail(),
		TitleID:        test.ValidTitleId,
		OgLanguageID:   test.ValidOgLanguageId,
		NewLanguageID:  test.ValidNewLanguageId,
		HashedPassword: hashedPassword,
	}
	return
}

func randomTranslate(phrase oapi.Phrase, languageId int16) db.Translate {

	return db.Translate{
		PhraseID:   phrase.Id,
		LanguageID: languageId,
		Phrase:     test.RandomString(8),
		PhraseHint: test.RandomString(8),
	}
}

func RandomTitle() (title db.Title) {

	return db.Title{
		ID:           test.RandomInt64(),
		Title:        test.RandomString(8),
		NumSubs:      test.RandomInt16(),
		OgLanguageID: validLanguageId,
	}
}

func randomLanguage() (language db.Language) {
	return db.Language{
		ID:       test.RandomInt16(),
		Language: test.RandomString(6),
		Tag:      "en",
	}
}

func setupHandlerTest(t *testing.T, ctrl *gomock.Controller, tc testCase, urlBasePath, body, method string) (*Server, echo.Context, *httptest.ResponseRecorder) {
	stubs := newMockStubs(ctrl)
	tc.buildStubs(stubs)

	e := echo.New()
	srv := NewServer(e, testCfg, stubs.mockQuerier, stubs.translateX, stubs.audioFileX)

	jwsToken, err := srv.fa.CreateJWSWithClaims(tc.permissions, tc.user)
	require.NoError(t, err)

	urlPath := urlBasePath + strconv.FormatInt(tc.extraInt, 10)

	req := handlerRequest(body, urlPath, method, string(jwsToken))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(token.UserIdContextKey, strconv.FormatInt(tc.user.ID, 10))

	return srv, c, rec
}

func setupServerTest(t *testing.T, ctrl *gomock.Controller, tc testCase) (*httptest.Server, string) {

	stubs := newMockStubs(ctrl)
	tc.buildStubs(stubs)

	e := echo.New()
	srv := NewServer(e, testCfg, stubs.mockQuerier, stubs.translateX, stubs.audioFileX)

	ts := httptest.NewServer(e)

	jwsToken, err := srv.fa.CreateJWSWithClaims(tc.permissions, tc.user)
	require.NoError(t, err)

	return ts, string(jwsToken)
}

func handlerRequest(json string, urlPath, method, authToken string) *http.Request {
	req := httptest.NewRequest(method, urlPath, strings.NewReader(json))

	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	req.Header.Set("Content-Type", "application/json-patch+json")

	return req
}

func jsonRequest(t *testing.T, json []byte, ts *httptest.Server, urlPath, method, authToken string) *http.Request {

	req, err := http.NewRequest(method, ts.URL+urlPath, bytes.NewBuffer(json))
	require.NoError(t, err)

	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	req.Header.Set("Content-Type", "application/json")

	return req
}

func createMultiPartBody(t *testing.T, data []byte, filename string, m map[string]string) (*bytes.Buffer, *multipart.Writer) {
	//data := []byte("This is the first sentence.\nThis is the second sentence.\n")
	err := os.WriteFile(filename, data, 0777)
	file, err := os.Open(filename)
	fmt.Println(file.Name())
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("filePath", filename)
	require.NoError(t, err)
	_, err = io.Copy(part, file)
	require.NoError(t, err)
	for key, val := range m {
		err = writer.WriteField(key, val)
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())
	return body, writer
}
