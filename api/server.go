package api

import (
	"encoding/json"
	mw "github.com/dsaker/nethttp-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-openapi/runtime/middleware"
	"log"
	"net/http"
	"sync"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/config"
	"talkliketv.click/tltv/internal/jsonlog"
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/token"
)

type Api struct {
	config  config.Config
	logger  *jsonlog.Logger
	queries db.Querier
	fa      *token.FakeAuthenticator
	Lock    sync.Mutex
}

// NewHandler creates a new HTTP server and sets up routing.
func NewHandler(cfg config.Config, logger *jsonlog.Logger, q db.Querier) (http.Handler, error) {

	r := http.NewServeMux()

	spec, err := oapi.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec\n: %s\n", err)
	}

	r.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		err = json.NewEncoder(w).Encode(spec)
		if err != nil {
			log.Fatalf("Error encoding swagger spec\n: %s\n", err)
		}
	})

	r.Handle("/swagger/", middleware.SwaggerUI(middleware.SwaggerUIOpts{
		Path:    "/swagger/",
		SpecURL: "/swagger/doc.json",
	}, nil))

	// Create a fake authenticator. This allows us to issue tokens, and also
	// implements a validator to check their validity.
	fa, err := token.NewFakeAuthenticator(&cfg.JWTDuration)
	if err != nil {
		log.Fatalln("error creating authenticator:", err)
	}

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	validator := mw.OapiRequestValidatorWithOptions(spec,
		&mw.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: token.NewAuthenticator(fa),
			},
		})

	h := oapi.HandlerWithOptions(
		&Api{
			config:  cfg,
			logger:  logger,
			queries: q,
			fa:      fa,
		},
		oapi.StdHTTPServerOptions{
			BaseURL:    "",
			BaseRouter: r,
			Middlewares: []oapi.MiddlewareFunc{
				validator,
			},
		},
	)

	return h, nil
}

// Make sure we conform to ServerInterface
var _ oapi.ServerInterface = (*Api)(nil)

// sendApiError wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendApiError(w http.ResponseWriter, code int, message string) {
	TitleErr := oapi.Error{
		Code:    int32(code),
		Message: message,
	}
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(TitleErr)
}

// sendApiError wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func (p *Api) sendInternalError(w http.ResponseWriter, err error) {
	p.logger.PrintError(err, nil)
	sendApiError(w, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func writeJSON(w http.ResponseWriter, status int, obj interface{}, headers http.Header) error {

	js, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}
