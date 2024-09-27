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

	// OgLanguageId Language id of title
	OgLanguageId int64 `json:"ogLanguageId"`

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

	// OgLanguageId Language id of title
	OgLanguageId int64 `json:"ogLanguageId"`

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
	Similarity string `form:"similarity" json:"similarity"`

	// Limit maximum number of results to return
	Limit int32 `form:"limit" json:"limit"`
}

// AddTitleJSONRequestBody defines body for AddTitle for application/json ContentType.
type AddTitleJSONRequestBody = NewTitle

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = NewUser

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
	CreateUser(ctx echo.Context) error
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
	// ------------- Required query parameter "similarity" -------------

	err = runtime.BindQueryParameter("form", true, true, "similarity", ctx.QueryParams(), &params.Similarity)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter similarity: %s", err))
	}

	// ------------- Required query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, true, "limit", ctx.QueryParams(), &params.Limit)
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

// CreateUser converts echo context to params.
func (w *ServerInterfaceWrapper) CreateUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateUser(ctx)
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
	router.POST(baseURL+"/users", wrapper.CreateUser)
	router.POST(baseURL+"/users/login", wrapper.LoginUser)
	router.DELETE(baseURL+"/users/:id", wrapper.DeleteUser)
	router.GET(baseURL+"/users/:id", wrapper.FindUserByID)
	router.PATCH(baseURL+"/users/:id", wrapper.UpdateUser)
	router.POST(baseURL+"/userspermissions", wrapper.AddUserPermission)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"AA==",
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
