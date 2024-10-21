//go:build go1.22

package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"talkliketv.click/tltv/api"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	"talkliketv.click/tltv/internal/jsonlog"
)

func main() {

	cfg := config.SetConfigs()

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// if ffmpeg is not installed and in PATH of host machine fail immediately
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	if !strings.Contains(string(output), "ffmpeg version") {
		logger.PrintFatal(errors.New(fmt.Sprintf("make sure ffmpeg is installed and in PATH: %s", string(output))), nil)
	}

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
