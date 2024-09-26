package api

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	v4mw "github.com/labstack/echo/v4/middleware"
	mw "github.com/oapi-codegen/echo-middleware"
	"log"
	"net/http"
	"sync"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	"talkliketv.click/tltv/internal/token"
)

type Server struct {
	sync.RWMutex
	queries db.Querier
	config  config.Config
	fa      token.FakeAuthenticator
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(e *echo.Echo, cfg config.Config, q db.Querier, spec *openapi3.T) *Server {

	// Create a fake authenticator. This allows us to issue tokens, and also
	// implements a validator to check their validity.
	fa, err := token.NewFakeAuthenticator(&cfg.JWTDuration)
	if err != nil {
		log.Fatalln("error creating authenticator:", err)
	}

	// Create middleware for validating tokens.
	middle, err := createMiddleware(fa, spec)
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}
	apiGrp := e.Group("/v1")
	apiGrp.Use(v4mw.Logger())
	apiGrp.Use(v4mw.Recover())
	apiGrp.Use(middle...)
	apiGrp.Use(v4mw.CORSWithConfig(v4mw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodOptions, http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	return &Server{
		fa:      *fa,
		queries: q,
		config:  cfg,
	}

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
