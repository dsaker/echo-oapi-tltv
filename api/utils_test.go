package api

import (
	"bytes"
	"flag"
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
	mockdb "talkliketv.click/tltv/db/mock"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	mockc "talkliketv.click/tltv/internal/mock/clients"
	mockt "talkliketv.click/tltv/internal/mock/translates"
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/test"
	"talkliketv.click/tltv/internal/token"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

var (
	testCfg            config.Config
	validLanguageModel = db.Language{
		ID:       27,
		Language: "English",
		Tag:      "en",
	}
)

const (
	usersBasePath           = "/v1/users"
	titlesBasePath          = "/v1/titles"
	usersPermissionBasePath = "/v1/userspermissions"
	phrasesBasePath         = "/v1/phrases"
	usersPhrasesBasePath    = "/v1/usersphrases"
	languagesBasePath       = "/v1/languages"
	validLanguageId         = 27
)

type testCase struct {
	name             string
	body             interface{}
	user             db.User
	userId           int64
	buildStubs       func(*mockdb.MockQuerier, *mockt.MockTranslateX)
	buildStubsClient func(*mockdb.MockQuerier, *mockt.MockTranslateX, *mockc.MockTranslateClientX)
	multipartBody    func(*testing.T) (*bytes.Buffer, *multipart.Writer)
	checkRecorder    func(rec *httptest.ResponseRecorder)
	checkResponse    func(res *http.Response)
	values           map[string]any
	permissions      []string
}

func TestMain(m *testing.M) {
	testCfg = config.SetConfigs()

	flag.Parse()
	testCfg.TTSBasePath = "/tmp/audio/"
	os.Exit(m.Run())
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

//func randomPhrase() oapi.Phrase {
//	return oapi.Phrase{
//		Id:      util.RandomInt64(),
//		TitleId: util.RandomInt64(),
//	}
//}

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
	text := mockt.NewMockTranslateX(ctrl)
	store := mockdb.NewMockQuerier(ctrl)
	tc.buildStubs(store, text)

	e := echo.New()
	srv := NewServer(e, testCfg, store, text)

	jwsToken, err := srv.fa.CreateJWSWithClaims(tc.permissions, tc.user)
	require.NoError(t, err)

	urlPath := urlBasePath + strconv.FormatInt(tc.userId, 10)

	req := handlerRequest(body, urlPath, method, string(jwsToken))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(token.UserIdContextKey, strconv.FormatInt(tc.user.ID, 10))

	return srv, c, rec
}

func setupServerTest(t *testing.T, ctrl *gomock.Controller, tc testCase) (*httptest.Server, string) {
	text := mockt.NewMockTranslateX(ctrl)
	store := mockdb.NewMockQuerier(ctrl)
	tc.buildStubs(store, text)

	e := echo.New()
	srv := NewServer(e, testCfg, store, text)

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
