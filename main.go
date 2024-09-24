//go:build go1.22

package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"github.com/getkin/kin-openapi/openapi3"
	ui "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo/v4"
	"log"
	"net"
	"net/http"
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

	e := echo.New()

	svr := api.NewServer(e, cfg, q)

	api.RegisterHandlers(e, svr)

	// Route
	e.GET("/hello", hello)

	e.Logger.Fatal(e.Start(net.JoinHostPort("0.0.0.0", cfg.Port)))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func runSwaggerUI(spec *openapi3.T) {
	r := http.NewServeMux()

	r.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(spec)
		if err != nil {
			log.Fatalf("Error encoding swagger spec\n: %s\n", err)
		}
	})

	r.Handle("/swagger/", ui.SwaggerUI(ui.SwaggerUIOpts{
		Path:    "/swagger/",
		SpecURL: "/swagger/doc.json",
	}, nil))

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", "8081"),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
