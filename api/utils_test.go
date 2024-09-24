package api

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"slices"
	"strings"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	"testing"
)

var (
	cfg config.Config
)

func TestMain(m *testing.M) {
	cfg = config.SetConfigs()

	flag.Parse()

	os.Exit(m.Run())
}

func newTestServer(e *echo.Echo, t *testing.T, q db.Querier) *Server {

	server := NewServer(e, cfg, q)

	return server
}

//func readBody(t *testing.T, rs *httptest.ResponseRecorder) string {
//	// Read the response body from the test server.
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			t.Fatal(err)
//		}
//	}(rs.Body)
//
//	body, err := io.ReadAll(rs.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//	bytes.TrimSpace(body)
//
//	return string(body)
//}

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
