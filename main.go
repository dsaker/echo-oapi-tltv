//go:build go1.22

package main

import (
	"database/sql"
	"flag"
	"net"
	"os"
	"talkliketv.click/tltv/api"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	"talkliketv.click/tltv/internal/jsonlog"
)

func main() {

	cfg := config.SetConfigs()

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

	e, _ := api.NewServer(cfg, q, &api.Translate{})

	e.Logger.Fatal(e.Start(net.JoinHostPort("0.0.0.0", cfg.Port)))
}
