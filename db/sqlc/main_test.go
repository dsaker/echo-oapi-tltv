package db

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"talkliketv.click/tltv/internal/util"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://tltvtest:pa55word@localhost/tltv_testdb?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(conn)
	flag.BoolVar(&util.Integration, "integration", false, "Run integration tests")
	flag.Parse()
	os.Exit(m.Run())
}
