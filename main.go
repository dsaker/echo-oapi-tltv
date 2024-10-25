//go:build go1.22

package main

import (
	"database/sql"
	"flag"
	"github.com/labstack/echo/v4"
	"net"
	"os/exec"
	"strings"
	"talkliketv.click/tltv/api"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	"talkliketv.click/tltv/internal/translates"
)

func main() {

	e := echo.New()

	cfg := config.SetConfigs()

	flag.Parse()

	// if ffmpeg is not installed and in PATH of host machine fail immediately
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		e.Logger.Fatal("Please make sure ffmep is installed and in PATH\n: %s", err)
	}
	if !strings.Contains(string(output), "ffmpeg version") {
		e.Logger.Fatal("Please make sure ffmep is installed and in PATH\n: %s", string(output))
	}

	// open db connection. if err fail immediately
	conn, err := cfg.OpenDB()
	if err != nil {
		e.Logger.Fatal("Error connecting to DB\n: %s", err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			e.Logger.Fatal("Error connecting to DB\n: %s", err)
		}
	}(conn)

	e.Logger.Info("database connection pool established")

	q := db.New(conn)

	api.NewServer(e, cfg, q, &translates.Translate{})

	e.Logger.Fatal(e.Start(net.JoinHostPort("0.0.0.0", cfg.Port)))
}
