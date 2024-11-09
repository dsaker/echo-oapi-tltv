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
	mockdb "talkliketv.click/tltv/internal/mock/db"
	mockt "talkliketv.click/tltv/internal/mock/translates"
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/test"
	"talkliketv.click/tltv/internal/token"
	"talkliketv.click/tltv/internal/util"
	"testing"
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
	voicesBasePath          = "/v1/voices"
)

type MockStubs struct {
	MockQuerier      *mockdb.MockQuerier
	TranslateX       *mockt.MockTranslateX
	TranslateClientX *mockt.MockTranslateClientX
	TtsClientX       *mockt.MockTTSClientX
	AudioFileX       *mocka.MockAudioFileX
}

// NewMockStubs creates instantiates new instances of all the mock interfaces for testing
func NewMockStubs(ctrl *gomock.Controller) MockStubs {
	return MockStubs{
		MockQuerier:      mockdb.NewMockQuerier(ctrl),
		TranslateX:       mockt.NewMockTranslateX(ctrl),
		TranslateClientX: mockt.NewMockTranslateClientX(ctrl),
		TtsClientX:       mockt.NewMockTTSClientX(ctrl),
		AudioFileX:       mocka.NewMockAudioFileX(ctrl),
	}
}

// testCase struct groups together the fields necessary for running most of the test
// cases
type testCase struct {
	name          string
	body          interface{}
	user          db.User
	extraInt      int64
	buildStubs    func(stubs MockStubs)
	multipartBody func(t *testing.T) (*bytes.Buffer, *multipart.Writer)
	checkRecorder func(rec *httptest.ResponseRecorder)
	checkResponse func(res *http.Response)
	values        map[string]any
	permissions   []string
	cleanUp       func(*testing.T)
}

func TestMain(m *testing.M) {
	_ = config.SetConfigs(&testCfg)
	flag.Parse()
	testCfg.TTSBasePath = test.AudioBasePath
	os.Exit(m.Run())
}

// readBody reads the http response body and returns it as a string
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

// randomUser creates a random user for testing
func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:             util.RandomInt64(),
		Name:           util.RandomString(8),
		Email:          util.RandomEmail(),
		TitleID:        util.ValidTitleId,
		OgLanguageID:   util.ValidOgLanguageId,
		NewLanguageID:  util.ValidNewLanguageId,
		HashedPassword: hashedPassword,
	}
	return
}

// randomTranslate create a random db Translate for testing
func randomTranslate(phrase oapi.Phrase, languageId int16) db.Translate {

	return db.Translate{
		PhraseID:   phrase.Id,
		LanguageID: languageId,
		Phrase:     util.RandomString(8),
		PhraseHint: util.RandomString(8),
	}
}

// randomLanguage creates a random db Language for testing
func randomLanguage() (language db.Language) {
	return db.Language{
		ID:       util.RandomInt16(),
		Language: util.RandomString(6),
		Tag:      "en",
	}
}

// randomVoice creates a random db Voice for testing
func randomVoice() (voice db.Voice) {
	return db.Voice{
		ID:                     util.RandomInt16(),
		LanguageID:             util.RandomInt16(),
		LanguageCodes:          []string{util.RandomString(8), util.RandomString(8)},
		SsmlGender:             "FEMALE",
		Name:                   util.RandomString(8),
		NaturalSampleRateHertz: 24000,
	}
}

// setupHandlerTest sets up a testCase that will be run through the handler
// these tests will not include the middleware JWT verification or the automated validation
// through openapi
func setupHandlerTest(t *testing.T, ctrl *gomock.Controller, tc testCase, urlBasePath, body, method string) (*Server, echo.Context, *httptest.ResponseRecorder) {
	stubs := NewMockStubs(ctrl)
	tc.buildStubs(stubs)

	e := echo.New()
	srv := NewServer(e, testCfg, stubs.MockQuerier, stubs.TranslateX, stubs.AudioFileX)

	jwsToken, err := srv.fa.CreateJWSWithClaims(tc.permissions, tc.user)
	require.NoError(t, err)

	urlPath := urlBasePath + strconv.FormatInt(tc.extraInt, 10)

	req := handlerRequest(body, urlPath, method, string(jwsToken))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(token.UserIdContextKey, strconv.FormatInt(tc.user.ID, 10))

	return srv, c, rec
}

// setupServerTest sets up testCase that will include the middleware not included in handler tests
func setupServerTest(t *testing.T, ctrl *gomock.Controller, tc testCase) (*httptest.Server, string) {

	stubs := NewMockStubs(ctrl)
	tc.buildStubs(stubs)

	e := echo.New()
	srv := NewServer(e, testCfg, stubs.MockQuerier, stubs.TranslateX, stubs.AudioFileX)

	ts := httptest.NewServer(e)

	jwsToken, err := srv.fa.CreateJWSWithClaims(tc.permissions, tc.user)
	require.NoError(t, err)

	return ts, string(jwsToken)
}

// handlerRequest is a helper function for setupHandlerTest
func handlerRequest(json string, urlPath, method, authToken string) *http.Request {
	req := httptest.NewRequest(method, urlPath, strings.NewReader(json))

	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	req.Header.Set("Content-Type", "application/json-patch+json")

	return req
}

// jsonRequest creates a new request which has json as the body and sets the Header content type to
// application/json
func jsonRequest(t *testing.T, json []byte, ts *httptest.Server, urlPath, method, authToken string) *http.Request {

	req, err := http.NewRequest(method, ts.URL+urlPath, bytes.NewBuffer(json))
	require.NoError(t, err)

	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	req.Header.Set("Content-Type", "application/json")

	return req
}

// createMultiPartBody creates and returns a multipart Writer.
// data is the data you want to write to the file.
// m is the map[string][string] of the fields, values you want to write to the multipart body
func createMultiPartBody(t *testing.T, data []byte, filename string, m map[string]string) (*bytes.Buffer, *multipart.Writer) {
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
