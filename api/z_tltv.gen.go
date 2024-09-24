// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for JSONPatchRequestAddReplaceTestOp.
const (
	Add     JSONPatchRequestAddReplaceTestOp = "add"
	Replace JSONPatchRequestAddReplaceTestOp = "replace"
	Test    JSONPatchRequestAddReplaceTestOp = "test"
)

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error message
	Message string `json:"message"`
}

// JSONPatchRequestAddReplaceTest defines model for JSONPatchRequestAddReplaceTest.
type JSONPatchRequestAddReplaceTest struct {
	// Op The operation to perform.
	Op JSONPatchRequestAddReplaceTestOp `json:"op"`

	// Path A JSON Pointer path.
	Path string `json:"path"`

	// Value The value to add, replace or test.
	Value interface{} `json:"value"`
}

// JSONPatchRequestAddReplaceTestOp The operation to perform.
type JSONPatchRequestAddReplaceTestOp string

// NewTitle defines model for NewTitle.
type NewTitle struct {
	// LanguageId Language id of title
	LanguageId int64 `json:"languageId"`

	// NumSubs Number of phrases
	NumSubs int32 `json:"numSubs"`

	// Title Name of the title
	Title string `json:"title"`
}

// NewUser defines model for NewUser.
type NewUser struct {
	// Email Email of user
	Email string `json:"email"`

	// Flipped switch between learning from or to native language
	Flipped bool `json:"flipped"`

	// Name Username of user. Must be alphanumeric.
	Name string `json:"name"`

	// NewLanguageId Id of language to learn
	NewLanguageId int64 `json:"newLanguageId"`

	// OgLanguageId Id of native language
	OgLanguageId int64 `json:"ogLanguageId"`

	// Password Password of user
	Password string `json:"password"`

	// TitleId Id of title to learn
	TitleId int64 `json:"titleId"`
}

// NewUserPermission defines model for NewUserPermission.
type NewUserPermission struct {
	// PermissionId Permission id of permission
	PermissionId int64 `json:"permission_id"`

	// UserId User id of user
	UserId int64 `json:"user_id"`
}

// PatchRequest defines model for PatchRequest.
type PatchRequest = []PatchRequest_Item

// PatchRequest_Item defines model for PatchRequest.Item.
type PatchRequest_Item struct {
	union json.RawMessage
}

// Title defines model for Title.
type Title struct {
	// Id Unique id of the title
	Id int64 `json:"id"`

	// LanguageId Language id of title
	LanguageId int64 `json:"languageId"`

	// NumSubs Number of phrases
	NumSubs int32 `json:"numSubs"`

	// Title Name of the title
	Title string `json:"title"`
}

// User defines model for User.
type User struct {
	// Email Email of user
	Email string `json:"email"`

	// Flipped switch between learning from or to native language
	Flipped bool `json:"flipped"`

	// Id Unique id of the user
	Id int64 `json:"id"`

	// Name Username of user. Must be alphanumeric.
	Name string `json:"name"`

	// NewLanguageId Id of language to learn
	NewLanguageId int64 `json:"newLanguageId"`

	// OgLanguageId Id of native language
	OgLanguageId int64 `json:"ogLanguageId"`

	// Password Password of user
	Password string `json:"password"`

	// TitleId Id of title to learn
	TitleId int64 `json:"titleId"`
}

// UserLogin defines model for UserLogin.
type UserLogin struct {
	// Password Password of user
	Password string `json:"password"`

	// Username Username of user
	Username string `json:"username"`
}

// UserLoginResponse defines model for UserLoginResponse.
type UserLoginResponse struct {
	// Jwt token of user
	Jwt string `json:"jwt"`
}

// UserPermissionResponse defines model for UserPermissionResponse.
type UserPermissionResponse struct {
	// Id Unique id of the user permission
	Id int64 `json:"id"`

	// PermissionId Permission id of permission
	PermissionId int64 `json:"permission_id"`

	// UserId User id of user
	UserId int64 `json:"user_id"`
}

// FindTitlesParams defines parameters for FindTitles.
type FindTitlesParams struct {
	// Similarity find titles similar to
	Similarity *string `form:"similarity,omitempty" json:"similarity,omitempty"`

	// Limit maximum number of results to return
	Limit *int32 `form:"limit,omitempty" json:"limit,omitempty"`
}

// AddTitleJSONRequestBody defines body for AddTitle for application/json ContentType.
type AddTitleJSONRequestBody = NewTitle

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody = NewUser

// LoginUserJSONRequestBody defines body for LoginUser for application/json ContentType.
type LoginUserJSONRequestBody = UserLogin

// UpdateUserApplicationJSONPatchPlusJSONRequestBody defines body for UpdateUser for application/json-patch+json ContentType.
type UpdateUserApplicationJSONPatchPlusJSONRequestBody = PatchRequest

// AddUserPermissionJSONRequestBody defines body for AddUserPermission for application/json ContentType.
type AddUserPermissionJSONRequestBody = NewUserPermission

// AsJSONPatchRequestAddReplaceTest returns the union data inside the PatchRequest_Item as a JSONPatchRequestAddReplaceTest
func (t PatchRequest_Item) AsJSONPatchRequestAddReplaceTest() (JSONPatchRequestAddReplaceTest, error) {
	var body JSONPatchRequestAddReplaceTest
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromJSONPatchRequestAddReplaceTest overwrites any union data inside the PatchRequest_Item as the provided JSONPatchRequestAddReplaceTest
func (t *PatchRequest_Item) FromJSONPatchRequestAddReplaceTest(v JSONPatchRequestAddReplaceTest) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeJSONPatchRequestAddReplaceTest performs a merge with any union data inside the PatchRequest_Item, using the provided JSONPatchRequestAddReplaceTest
func (t *PatchRequest_Item) MergeJSONPatchRequestAddReplaceTest(v JSONPatchRequestAddReplaceTest) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t PatchRequest_Item) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *PatchRequest_Item) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Returns all titles
	// (GET /titles)
	FindTitles(ctx echo.Context, params FindTitlesParams) error
	// Creates a new title
	// (POST /titles)
	AddTitle(ctx echo.Context) error
	// Deletes a title by ID
	// (DELETE /titles/{id})
	DeleteTitle(ctx echo.Context, id int64) error
	// Returns a title by ID
	// (GET /titles/{id})
	FindTitleByID(ctx echo.Context, id int64) error
	// Creates a new user
	// (POST /users)
	Register(ctx echo.Context) error
	// Login a user
	// (POST /users/login)
	LoginUser(ctx echo.Context) error
	// Deletes a user by ID
	// (DELETE /users/{id})
	DeleteUser(ctx echo.Context, id int64) error
	// Returns a user by ID
	// (GET /users/{id})
	FindUserByID(ctx echo.Context, id int64) error
	// Patch an existing user
	// (PATCH /users/{id})
	UpdateUser(ctx echo.Context, id int64) error

	// (POST /userspermissions)
	AddUserPermission(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// FindTitles converts echo context to params.
func (w *ServerInterfaceWrapper) FindTitles(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"titles:r"})

	// Parameter object where we will unmarshal all parameters from the context
	var params FindTitlesParams
	// ------------- Optional query parameter "similarity" -------------

	err = runtime.BindQueryParameter("form", true, false, "similarity", ctx.QueryParams(), &params.Similarity)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter similarity: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.FindTitles(ctx, params)
	return err
}

// AddTitle converts echo context to params.
func (w *ServerInterfaceWrapper) AddTitle(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"titles:w"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AddTitle(ctx)
	return err
}

// DeleteTitle converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteTitle(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteTitle(ctx, id)
	return err
}

// FindTitleByID converts echo context to params.
func (w *ServerInterfaceWrapper) FindTitleByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.FindTitleByID(ctx, id)
	return err
}

// Register converts echo context to params.
func (w *ServerInterfaceWrapper) Register(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Register(ctx)
	return err
}

// LoginUser converts echo context to params.
func (w *ServerInterfaceWrapper) LoginUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.LoginUser(ctx)
	return err
}

// DeleteUser converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteUser(ctx, id)
	return err
}

// FindUserByID converts echo context to params.
func (w *ServerInterfaceWrapper) FindUserByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.FindUserByID(ctx, id)
	return err
}

// UpdateUser converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateUser(ctx, id)
	return err
}

// AddUserPermission converts echo context to params.
func (w *ServerInterfaceWrapper) AddUserPermission(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"global:admin"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AddUserPermission(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/titles", wrapper.FindTitles)
	router.POST(baseURL+"/titles", wrapper.AddTitle)
	router.DELETE(baseURL+"/titles/:id", wrapper.DeleteTitle)
	router.GET(baseURL+"/titles/:id", wrapper.FindTitleByID)
	router.POST(baseURL+"/users", wrapper.Register)
	router.POST(baseURL+"/users/login", wrapper.LoginUser)
	router.DELETE(baseURL+"/users/:id", wrapper.DeleteUser)
	router.GET(baseURL+"/users/:id", wrapper.FindUserByID)
	router.PATCH(baseURL+"/users/:id", wrapper.UpdateUser)
	router.POST(baseURL+"/userspermissions", wrapper.AddUserPermission)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"OLJxqBhYVLg/FzbOxxn/CQAA///8jBu+fiYAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
