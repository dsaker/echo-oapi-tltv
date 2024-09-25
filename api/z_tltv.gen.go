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
