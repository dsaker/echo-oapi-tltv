package api

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	ui "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo/v4"
	v4mw "github.com/labstack/echo/v4/middleware"
	mw "github.com/oapi-codegen/echo-middleware"
	"log"
	"sync"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	"talkliketv.click/tltv/internal/token"
)

type Server struct {
	sync.RWMutex
	queries    db.Querier
	Translates TranslateX
	config     config.Config
	fa         token.FakeAuthenticator
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(e *echo.Echo, cfg config.Config, q db.Querier, t TranslateX) *Server {

	spec, err := GetSwagger()
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
	fa, err := token.NewFakeAuthenticator(&cfg.JWTDuration)
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

	srv := &Server{
		fa:         *fa,
		Translates: t,
		queries:    q,
		config:     cfg,
	}

	RegisterHandlersWithBaseURL(apiGrp, srv, "")

	return srv
}

// Make sure we conform to ServerInterface
var _ ServerInterface = (*Server)(nil)

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
