//go:build go1.22

package main

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/translate"
	"context"
	"database/sql"
	"flag"
	"talkliketv.click/tltv/internal/translates"

	//"github.com/kataras/iris/v12/context"
	"github.com/labstack/echo/v4"
	"net"
	"os/exec"
	"strings"
	"talkliketv.click/tltv/api"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/audio/audiofile"
	"talkliketv.click/tltv/internal/config"
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
	ctx := context.Background()
	// create translate client
	transClient, err := translate.NewClient(ctx)
	if err != nil {
		e.Logger.Fatal("Error creating google api translate client\n: %s", err)
	}
	ttsClient, err := texttospeech.NewClient(ctx)
	if err != nil {
		e.Logger.Fatal("Error creating google api translate client\n: %s", err)
	}
	t := translates.New(transClient, ttsClient)
	api.NewServer(e, cfg, q, t, &audiofile.AudioFile{})

	e.Logger.Fatal(e.Start(net.JoinHostPort("0.0.0.0", cfg.Port)))
}
