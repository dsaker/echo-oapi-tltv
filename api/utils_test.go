package api

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"slices"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	"talkliketv.click/tltv/internal/jsonlog"
	"talkliketv.click/tltv/internal/token"
	"testing"
	"time"
)

var (
	cfg    config.Config
	logger *jsonlog.Logger
)

func TestMain(m *testing.M) {
	cfg = config.SetConfigs()
	// get port and debug from commandline flags... if not present use defaults
	flag.IntVar(&cfg.Port, "port", 8080, "API server port")

	flag.Parse()

	logger = jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	os.Exit(m.Run())
}

func newTestServer(t *testing.T, q db.Querier) (*httptest.Server, *token.FakeAuthenticator) {

	h, err := NewHandler(cfg, logger, q)
	require.NoError(t, err)

	ts := httptest.NewServer(h)

	duration := time.Hour * 24
	fa, err := token.NewFakeAuthenticator(&duration)
	if err != nil {
		log.Fatalln("error creating authenticator:", err)
	}

	return ts, fa
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

func request(t *testing.T, json []byte, ts *httptest.Server, urlPath, method, authToken string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+urlPath, bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	req.Header.Set("Content-Type", "application/json-patch+json")
	res, err := ts.Client().Do(req)
	if err != nil {
		fmt.Printf("response error: %s: %s\n", err, readBody(t, res))
	}

	return res
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
