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
	"reflect"
	"slices"
	"strconv"
	"strings"
	mockdb "talkliketv.click/tltv/db/mock"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	mock "talkliketv.click/tltv/internal/mock"
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
	validLanguageId         = 27
)

type testCase struct {
	name          string
	body          interface{}
	user          db.User
	userId        int64
	buildStubs    func(*mockdb.MockQuerier, *mock.MockTranslateX)
	multipartBody func(*testing.T) (*bytes.Buffer, *multipart.Writer)
	checkRecorder func(rec *httptest.ResponseRecorder)
	checkResponse func(res *http.Response)
	values        map[string]any
	permissions   []string
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

func randomPhrase() Phrase {
	return Phrase{
		Id:      util.RandomInt64(),
		TitleId: util.RandomInt64(),
	}
}

func randomTranslate(phrase Phrase, languageId int16) db.Translate {

	return db.Translate{
		PhraseID:   phrase.Id,
		LanguageID: languageId,
		Phrase:     util.RandomString(8),
		PhraseHint: util.RandomString(8),
	}
}

func randomTitle() (title db.Title) {

	return db.Title{
		ID:           util.RandomInt64(),
		Title:        util.RandomString(8),
		NumSubs:      util.RandomInt16(),
		OgLanguageID: validLanguageId,
	}
}

func randomLanguage() (language db.Language) {
	return db.Language{
		ID:       util.RandomInt16(),
		Language: util.RandomString(6),
		Tag:      "en",
	}
}

func setupHandlerTest(t *testing.T, ctrl *gomock.Controller, tc testCase, urlBasePath, body, method string) (*Server, echo.Context, *httptest.ResponseRecorder) {
	text := mock.NewMockTranslateX(ctrl)
	store := mockdb.NewMockQuerier(ctrl)
	tc.buildStubs(store, text)

	e, srv := NewServer(testCfg, store, text)

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
	text := mock.NewMockTranslateX(ctrl)
	store := mockdb.NewMockQuerier(ctrl)
	tc.buildStubs(store, text)

	e, srv := NewServer(testCfg, store, text)

	ts := httptest.NewServer(e)

	jwsToken, err := srv.fa.CreateJWSWithClaims(tc.permissions, tc.user)
	require.NoError(t, err)

	//req := jsonRequest(t, body, ts, urlPath, method, string(jwsToken))
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

func requireMatchAnyExcept(t *testing.T, model any, response any, skip []string, except string, shouldEqual any) {

	v := reflect.ValueOf(response)
	u := reflect.ValueOf(model)

	for i := 0; i < v.NumField(); i++ {
		// Check if field name is the one that should be different
		if v.Type().Field(i).Name == except {
			// Check if type is int32 or int64
			if v.Field(i).CanInt() {
				// check if equal as int64
				require.Equal(t, shouldEqual, v.Field(i).Int())
			} else {
				// if not check if equal as string
				require.Equal(t, shouldEqual, v.Field(i).String())
			}
		} else if slices.Contains(skip, v.Type().Field(i).Name) {
			continue
		} else {
			if v.Field(i).CanInt() {
				require.Equal(t, u.Field(i).Int(), v.Field(i).Int())
			} else {
				require.Equal(t, u.Field(i).String(), v.Field(i).String())
			}
		}
	}
}
