package api

import (
	"bytes"
	"flag"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"io"
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
	"talkliketv.click/tltv/internal/token"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

var (
	testCfg config.Config
)

func TestMain(m *testing.M) {
	testCfg = config.SetConfigs()

	flag.Parse()

	os.Exit(m.Run())
}

func readBody(t *testing.T, rs *http.Response) string {
	// Read the response body from the test server.
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
		ID:             util.RandomInt64(1, 1000),
		Name:           util.RandomName(),
		Email:          util.RandomEmail(),
		TitleID:        util.ValidTitleId,
		Flipped:        false,
		OgLanguageID:   util.ValidOgLanguageId,
		NewLanguageID:  util.ValidNewLanguageId,
		HashedPassword: hashedPassword,
	}
	return
}

func setupTest(t *testing.T, ctrl *gomock.Controller, tc testCase, body, method string) (*Server, echo.Context, *httptest.ResponseRecorder) {
	store := mockdb.NewMockQuerier(ctrl)
	tc.buildStubs(store)

	spec, err := GetSwagger()
	require.NoError(t, err)

	e := echo.New()
	srv := NewServer(e, testCfg, store, spec)

	jwsToken, err := srv.fa.CreateJWSWithClaims([]string{db.ReadTitlesCode}, tc.user)
	require.NoError(t, err)

	urlPath := util.UserBasePath + strconv.FormatInt(tc.userId, 10)

	req := request(body, urlPath, method, string(jwsToken))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(token.UserIdContextKey, strconv.FormatInt(tc.user.ID, 10))

	return srv, c, rec
}

func request(json string, urlPath, method, authToken string) *http.Request {
	req := httptest.NewRequest(method, urlPath, strings.NewReader(json))
	//req := httptest.

	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	req.Header.Set("Content-Type", "application/json-patch+json")

	return req
}

func requireMatchAnyExcept(t *testing.T, model any, response any, skip []string, except, shouldEqual string) {

	v := reflect.ValueOf(response)
	u := reflect.ValueOf(model)
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Name == except {
			require.Equal(t, shouldEqual, v.Field(i).String())
		} else if slices.Contains(skip, v.Type().Field(i).Name) {
			continue
		} else {
			require.Equal(t, u.Field(i).String(), v.Field(i).String())
		}
	}
}
