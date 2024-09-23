//go:build go1.22

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"talkliketv.click/tltv/api"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	"talkliketv.click/tltv/internal/jsonlog"
)

func main() {

	cfg := config.SetConfigs()
	// get port and debug from commandline flags... if not present use defaults
	flag.IntVar(&cfg.Port, "port", 8080, "API server port")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// open db connection. if err fail immediately
	conn, err := cfg.OpenDB()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.PrintFatal(err, nil)
		}
	}(conn)

	logger.PrintInfo("database connection pool established", nil)

	q := db.New(conn)

	handler, err := api.NewHandler(cfg, logger, q)

	s := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", cfg.Port),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
