package api

//import (
//	"encoding/json"
//	"fmt"
//	middleware "github.com/oapi-codegen/nethttp-middleware"
//	"github.com/oapi-codegen/testutil"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"net/http"
//	"net/http/httptest"
//	"talkliketv.click/tltv/internal/oapi"
//	"testing"
//)
//
//func doGet(t *testing.T, mux *http.ServeMux, url string) *httptest.ResponseRecorder {
//	response := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, mux)
//	return response.Recorder
//}
//
//func TestTlTv(t *testing.T) {
//	var err error
//
//	// Get the swagger description of our API
//	swagger, err := oapi.GetSwagger()
//	require.NoError(t, err)
//
//	// Clear out the servers array in the swagger spec, that skips validating
//	// that server names match. We don't know how this thing will be run.
//	swagger.Servers = nil
//
//	// Create a new ServeMux for testing.
//	m := http.NewServeMux()
//
//	// Use our validation middleware to check all requests against the
//	// OpenAPI schema.
//	opts := oapi.StdHTTPServerOptions{
//		BaseRouter: m,
//		Middlewares: []oapi.MiddlewareFunc{
//			middleware.OapiRequestValidator(swagger),
//		},
//	}
//
//	app := newTest Application(t)
//	oapi.HandlerWithOptions(app, opts)
//
//	t.Run("Add title", func(t *testing.T) {
//		newTitle := oapi.NewTitle{
//			LanguageId: 1,
//			NumSubs:    1,
//			Title:      "new title",
//		}
//
//		rr := testutil.NewRequest().Post("/titles").WithJsonBody(newTitle).GoWithHTTPHandler(t, m).Recorder
//		assert.Equal(t, http.StatusCreated, rr.Code)
//
//		var resultTitle oapi.Title
//		err = json.NewDecoder(rr.Body).Decode(&resultTitle)
//		assert.NoError(t, err, "error unmarshalling response")
//		assert.Equal(t, newTitle.Title, resultTitle.Title)
//		assert.Equal(t, newTitle.NumSubs, resultTitle.NumSubs)
//		assert.Equal(t, newTitle.LanguageId, resultTitle.LanguageId)
//	})
//
//	t.Run("Find title by ID", func(t *testing.T) {
//		title := oapi.Title{
//			Id: 100,
//		}
//
//		Server.Titles[title.Id] = title
//		rr := doGet(t, m, fmt.Sprintf("/titles/%d", title.Id))
//
//		var resultTitle oapi.Title
//		err = json.NewDecoder(rr.Body).Decode(&resultTitle)
//		assert.NoError(t, err, "error getting title")
//		assert.Equal(t, title, resultTitle)
//	})
//
//	t.Run("Title not found", func(t *testing.T) {
//		rr := doGet(t, m, "/titles/27179095781")
//		assert.Equal(t, http.StatusNotFound, rr.Code)
//
//		var titleError oapi.Error
//		err = json.NewDecoder(rr.Body).Decode(&titleError)
//		assert.NoError(t, err, "error getting response", err)
//		assert.Equal(t, int32(http.StatusNotFound), titleError.Code)
//	})
//
//	t.Run("List all titles", func(t *testing.T) {
//		Server.Titles = map[int64]oapi.Title{
//			1: {},
//			2: {},
//		}
//
//		// Now, list all titles, we should have two
//		rr := doGet(t, m, "/titles")
//		assert.Equal(t, http.StatusOK, rr.Code)
//
//		var titleList []oapi.Title
//		err = json.NewDecoder(rr.Body).Decode(&titleList)
//		assert.NoError(t, err, "error getting response", err)
//		assert.Equal(t, 2, len(titleList))
//	})
//
//	t.Run("Delete titles", func(t *testing.T) {
//		Server.Titles = map[int64]oapi.Title{
//			1: {},
//			2: {},
//		}
//
//		// Let's delete non-existent title
//		rr := testutil.NewRequest().Delete("/titles/7").GoWithHTTPHandler(t, m).Recorder
//		assert.Equal(t, http.StatusNotFound, rr.Code)
//
//		var titleError oapi.Error
//		err = json.NewDecoder(rr.Body).Decode(&titleError)
//		assert.NoError(t, err, "error unmarshalling TitleError")
//		assert.Equal(t, int32(http.StatusNotFound), titleError.Code)
//
//		// Now, delete both real titles
//		rr = testutil.NewRequest().Delete("/titles/1").GoWithHTTPHandler(t, m).Recorder
//		assert.Equal(t, http.StatusNoContent, rr.Code)
//
//		rr = testutil.NewRequest().Delete("/titles/2").GoWithHTTPHandler(t, m).Recorder
//		assert.Equal(t, http.StatusNoContent, rr.Code)
//
//		// Should have no titles left.
//		var titleList []oapi.Title
//		rr = doGet(t, m, "/titles")
//		assert.Equal(t, http.StatusOK, rr.Code)
//		err = json.NewDecoder(rr.Body).Decode(&titleList)
//		assert.NoError(t, err, "error getting response", err)
//		assert.Equal(t, 0, len(titleList))
//	})
//}
