package api

import (
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/mock/gomock"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	mocka "talkliketv.click/tltv/internal/mock/audiofile"
	mockdb "talkliketv.click/tltv/internal/mock/db"
	mockt "talkliketv.click/tltv/internal/mock/translates"
	"talkliketv.click/tltv/internal/test"
	"talkliketv.click/tltv/internal/token"
	"talkliketv.click/tltv/internal/util"
)

var (
	testCfg     TestConfig
	count       = 0
	mappedPort  nat.Port
	integration = false
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
	templateDb              = "templatedb"
	dbUser                  = "postgres"
	dbPass                  = "postgres"
	dbPort                  = "5432/tcp"
	pgData                  = "/var/lib/pg/data"
)

type MockStubs struct {
	MockQuerier      *mockdb.MockQuerier
	TranslateX       *mockt.MockTranslateX
	TranslateClientX *mockt.MockTranslateClientX
	TtsClientX       *mockt.MockTTSClientX
	AudioFileX       *mocka.MockAudioFileX
}

type TestConfig struct {
	config.Config
	conn      *sql.DB
	container testcontainers.Container
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
	_ = config.SetConfigs(&testCfg.Config)
	flag.BoolVar(&integration, "integration", false, "Run integration tests")
	flag.Parse()
	testCfg.TTSBasePath = test.AudioBasePath
	if integration {
		testCfg.container, testCfg.conn = setupTemplateDb()
	}
	// Run the tests
	exitCode := m.Run()
	if integration {
		err := testCfg.container.Terminate(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}
	os.Exit(exitCode)
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

// randomLanguage creates a random db Language for testing
func randomLanguage() (language db.Language) {
	return db.Language{
		ID:       test.RandomInt16(),
		Language: test.RandomString(6),
		Tag:      "en",
	}
}

// setupHandlerTest sets up a testCase that will be run through the handler
// these tests will not include the middleware JWT verification or the automated validation
// through openapi
func setupHandlerTest(t *testing.T, ctrl *gomock.Controller, tc testCase, urlBasePath, body, method string) (*Server, echo.Context, *httptest.ResponseRecorder) {
	stubs := NewMockStubs(ctrl)
	tc.buildStubs(stubs)

	e := echo.New()
	srv := NewServer(e, testCfg.Config, stubs.MockQuerier, stubs.TranslateX, stubs.AudioFileX)

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
	srv := NewServer(e, testCfg.Config, stubs.MockQuerier, stubs.TranslateX, stubs.AudioFileX)

	ts := httptest.NewServer(e)

	jwsToken, err := srv.fa.CreateJWSWithClaims(tc.permissions, tc.user)
	require.NoError(t, err)

	return ts, string(jwsToken)
}

// setupIntegrationTest
func setupIntegrationTest(t *testing.T, e *echo.Echo, s *Server, tc testCase) (*httptest.Server, string) {
	ts := httptest.NewServer(e)

	jwsToken, err := s.fa.CreateJWSWithClaims(tc.permissions, tc.user)
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
	err := os.WriteFile(filename, data, 0600)
	require.NoError(t, err)
	file, err := os.Open(filename)
	require.NoError(t, err)
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

func setupTemplateDb() (testcontainers.Container, *sql.DB) {
	// setup db container
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	container, conn := createDbContainer()

	// mark testdb as a template so you can copy it
	query := `update pg_database set datistemplate=true where datname=$1;`
	_, err := conn.ExecContext(ctx, query, templateDb)
	if err != nil {
		log.Fatal(err)
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPass, mappedPort.Port(), templateDb)

	// migrate template database
	m, err := migrate.New("file://../db/migrations", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("migration done")

	// close connection with template db
	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	// create and return connection with postgres database
	conn = createConnection("postgres")
	return container, conn
}

func createDbContainer() (testcontainers.Container, *sql.DB) {
	var env = map[string]string{
		"POSTGRES_PASSWORD": dbPass,
		"POSTGRES_USER":     dbUser,
		"POSTGRES_DB":       templateDb,
		"PGDATA":            pgData,
	}

	// create a database in memory for faster copying
	// https://gajus.com/blog/setting-up-postgre-sql-for-running-integration-tests
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14-alpine",
			ExposedPorts: []string{dbPort},
			Env:          env,
			Tmpfs:        map[string]string{pgData: "rw"},
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	mappedPort, err = container.MappedPort(ctx, dbPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("postgres container ready and running at port: ", mappedPort.Port())

	time.Sleep(time.Second)

	//postgres://postgres:postgres@localhost:5433/tltv_testdb?sslmode=disable
	conn := createConnection(templateDb)
	return container, conn
}

func createTestDb(t *testing.T, m *Server) (*sql.DB, string) {
	ctx := context.Background()

	query := "create database testdb template templatedb"
	// increase count so each test db has a different name
	count++
	// lock so you don't have multiple connections to template db
	m.Lock()
	defer m.Unlock()

	//newDbName := templateDb + strconv.Itoa(count)
	_, err := testCfg.conn.ExecContext(ctx, query)
	if err != nil {
		t.Fatal(err)
	}

	conn := createConnection("testdb")
	return conn, "testdb"
}

func destroyDb(t *testing.T, c *sql.DB, dbName string) {
	err := c.Close()
	if err != nil {
		t.Fatal(err)
	}
	query := `DROP DATABASE testdb`
	_, err = testCfg.conn.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

func createConnection(dbName string) *sql.DB {
	connString := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPass, mappedPort.Port(), dbName)
	log.Println("created connection string: ", connString)
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), testCfg.Config.CtxTimeout)
	defer cancel()
	err = conn.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
