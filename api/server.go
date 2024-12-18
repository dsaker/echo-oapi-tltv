package api

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"os"
	"sync"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/translate"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	ui "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo/v4"
	v4mw "github.com/labstack/echo/v4/middleware"
	mw "github.com/oapi-codegen/echo-middleware"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/audio"
	"talkliketv.click/tltv/internal/audio/audiofile"
	"talkliketv.click/tltv/internal/config"
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/token"
	"talkliketv.click/tltv/internal/translates"
	"talkliketv.click/tltv/internal/util"
)

type Server struct {
	sync.RWMutex
	queries    db.Querier
	translates translates.TranslateX
	config     config.Config
	fa         token.FakeAuthenticator
	af         audiofile.AudioFileX
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(e *echo.Echo, cfg config.Config, q db.Querier, t translates.TranslateX, af audiofile.AudioFileX) *Server {
	// make sure silence mp3s exist in your base path
	initSilence(e, cfg)

	spec, err := oapi.GetSwagger()
	if err != nil {
		log.Fatalln("loading spec: %w", err)
	}

	g := e.Group("/swagger")
	g.GET("/doc.json", func(ctx echo.Context) error {
		err := json.NewEncoder(ctx.Response().Writer).Encode(spec)
		if err != nil {
			log.Fatalf("Error encoding swagger spec\n: %s\n", err)
		}
		return nil
	})

	swaggerHandler := ui.SwaggerUI(ui.SwaggerUIOpts{
		Path:    "/swagger/",
		SpecURL: "/swagger/doc.json",
	}, nil)

	g.GET("/", echo.WrapHandler(swaggerHandler))

	// Create a fake authenticator. This allows us to issue tokens, and also
	// implements a validator to check their validity.
	keyData, err := os.ReadFile(cfg.PrivateKeyPath) // Replace "private.key" with your key file name
	if err != nil {
		log.Fatal("Error reading key file:", err)
	}

	fa, err := token.NewFakeAuthenticator(&cfg.JWTDuration, keyData)
	if err != nil {
		log.Fatalln("error creating authenticator:", err)
	}

	spec.Servers = openapi3.Servers{&openapi3.Server{URL: "/v1"}}

	// Create middleware for validating tokens.
	middle, err := createMiddleware(fa, spec)
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}
	apiGrp := e.Group("/v1")
	apiGrp.Use(v4mw.Logger())
	apiGrp.Use(v4mw.Recover())
	apiGrp.Use(middle...)

	// add rate limiting middleware
	apiGrp.Use(v4mw.RateLimiter(v4mw.NewRateLimiterMemoryStore(rate.Limit(20))))

	srv := &Server{
		fa:         *fa,
		translates: t,
		queries:    q,
		config:     cfg,
		af:         af,
	}

	oapi.RegisterHandlersWithBaseURL(apiGrp, srv, "")

	return srv
}

// Make sure we conform to ServerInterface
var _ oapi.ServerInterface = (*Server)(nil)

// createMiddleware creates the JWS middleware function that will validate the JWT token
// and store data in echo context for use in echo handlers
func createMiddleware(v token.JWSValidator, spec *openapi3.T) ([]echo.MiddlewareFunc, error) {
	validator := mw.OapiRequestValidatorWithOptions(spec,
		&mw.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: token.NewAuthenticator(v),
			},
			SilenceServersWarning: true,
		})

	return []echo.MiddlewareFunc{validator}, nil
}

// initSilence copies the silence mp3's from the embedded filesystem to the config TTSBasePath
func initSilence(e *echo.Echo, cfg config.Config) {
	// check if silence mp3s exist in your base path
	silencePath := cfg.TTSBasePath + audiofile.AudioPauseFilePath[cfg.PhrasePause]
	exists, err := util.PathExists(silencePath)
	if err != nil {
		e.Logger.Fatal(err)
	}
	// if it doesn't exist copy it from embedded FS to TTSBasePath
	if !exists {
		err = os.MkdirAll(cfg.TTSBasePath+"silence/", 0777)
		if err != nil {
			e.Logger.Fatal(err)
		}
		for key, value := range audiofile.AudioPauseFilePath {
			fmt.Printf("%d", key)
			pause, err := audio.Silence.ReadFile(value)
			if err != nil {
				e.Logger.Fatal(err)
			}
			// Create a new file
			file, err := os.Create(cfg.TTSBasePath + value)
			if err != nil {
				e.Logger.Fatal(err)
			}
			defer file.Close()
			// Write to the file
			_, err = file.Write(pause)
			if err != nil {
				e.Logger.Fatal(err)
			}
			// Ensure data is written to disk
			err = file.Sync()
			if err != nil {
				e.Logger.Fatal(err)
			}
		}
	}
}

// CreateDependencies creates a new google translate and text-to-speech clients; constructs
// the translates and audiofile dependencies and returns them
func CreateDependencies(e *echo.Echo) (*translates.Translate, *audiofile.AudioFile) {
	// create google translate and text-to-speech clients
	ctx := context.Background()
	transClient, err := translate.NewClient(ctx)
	if err != nil {
		e.Logger.Fatal("Error creating google api translate client\n: %s", err)
	}
	ttsClient, err := texttospeech.NewClient(ctx)
	if err != nil {
		e.Logger.Fatal("Error creating google api translate client\n: %s", err)
	}
	t := translates.New(transClient, ttsClient)

	//initialize audiofile with the real command runner
	af := audiofile.New(&audiofile.RealCmdRunner{})

	return t, af
}
