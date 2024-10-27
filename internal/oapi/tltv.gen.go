// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package oapi

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
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for PatchRequestAddReplaceTestOp.
const (
	Add     PatchRequestAddReplaceTestOp = "add"
	Replace PatchRequestAddReplaceTestOp = "replace"
	Test    PatchRequestAddReplaceTestOp = "test"
)

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error message
	Message string `json:"message"`
}

// Language defines model for Language.
type Language struct {
	// Id id of language
	Id int16 `json:"id"`

	// Language string of language
	Language string `json:"language"`

	// Tag language tag used for google language methods
	Tag string `json:"tag"`
}

// NewTitle defines model for NewTitle.
type NewTitle struct {
	// Filename the file to upload
	Filename openapi_types.File `json:"filename"`

	// OgLanguageId Language id of title
	OgLanguageId int16 `json:"ogLanguageId"`

	// Title Name of the title
	Title string `json:"title"`
}

// NewUser defines model for NewUser.
type NewUser struct {
	// Email Email of user
	Email string `json:"email"`

	// Name Username of user. Must be alphanumeric.
	Name string `json:"name"`

	// NewLanguageId Id of language to learn
	NewLanguageId int16 `json:"newLanguageId"`

	// OgLanguageId Id of native language
	OgLanguageId int16 `json:"ogLanguageId"`

	// Password Password of user
	Password string `json:"password"`

	// TitleId Id of title to learn
	TitleId int64 `json:"titleId"`
}

// NewUserPermission defines model for NewUserPermission.
type NewUserPermission struct {
	// PermissionId Permission id of permission
	PermissionId int16 `json:"permissionId"`

	// UserId User id of user
	UserId int64 `json:"userId"`
}

// PatchRequest defines model for PatchRequest.
type PatchRequest = []PatchRequest_Item

// PatchRequest_Item defines model for PatchRequest.Item.
type PatchRequest_Item struct {
	union json.RawMessage
}

// PatchRequestAddReplaceTest defines model for PatchRequestAddReplaceTest.
type PatchRequestAddReplaceTest struct {
	// Op The operation to perform.
	Op PatchRequestAddReplaceTestOp `json:"op"`

	// Path A JSON Pointer path.
	Path string `json:"path"`

	// Value The value to add, replace or test.
	Value interface{} `json:"value"`
}

// PatchRequestAddReplaceTestOp The operation to perform.
type PatchRequestAddReplaceTestOp string

// Phrase defines model for Phrase.
type Phrase struct {
	// Id id of phrase
	Id int64 `json:"id"`

	// TitleId id of movie
	TitleId int64 `json:"titleId"`
}

// Title defines model for Title.
type Title struct {
	// Filename the file to upload
	Filename openapi_types.File `json:"filename"`

	// Id Unique id of the title
	Id int64 `json:"id"`

	// OgLanguageId Language id of title
	OgLanguageId int16 `json:"ogLanguageId"`

	// Title Name of the title
	Title string `json:"title"`
}

// TitlesTranslateRequest defines model for TitlesTranslateRequest.
type TitlesTranslateRequest struct {
	// NewLanguageId id of language to translate to
	NewLanguageId int16 `json:"newLanguageId"`

	// OldLanguageId id of language to translate from
	OldLanguageId int16 `json:"oldLanguageId"`

	// TitleId title id of title to translate from
	TitleId int64 `json:"titleId"`
}

// Translates defines model for Translates.
type Translates struct {
	LanguageId int16  `json:"languageId"`
	Phrase     string `json:"phrase"`
	PhraseHint string `json:"phraseHint"`
	PhraseId   int64  `json:"phraseId"`
}

// User defines model for User.
type User struct {
	// Email Email of user
	Email string `json:"email"`

	// Id Unique id of the user
	Id int64 `json:"id"`

	// Name Username of user. Must be alphanumeric.
	Name string `json:"name"`

	// NewLanguageId Id of language to learn
	NewLanguageId int16 `json:"newLanguageId"`

	// OgLanguageId Id of native language
	OgLanguageId int16 `json:"ogLanguageId"`

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
	Id int16 `json:"id"`

	// PermissionId Permission id of permission
	PermissionId int16 `json:"permissionId"`

	// UserId User id of user
	UserId int64 `json:"userId"`
}

// UsersPhrases defines model for UsersPhrases.
type UsersPhrases struct {
	// LanguageId id of language
	LanguageId int16 `json:"languageId"`

	// PhraseCorrect id of language
	PhraseCorrect int16 `json:"phraseCorrect"`

	// PhraseId id of phrase
	PhraseId int64 `json:"phraseId"`

	// TitleId id of title
	TitleId int64 `json:"titleId"`

	// UserId id of user
	UserId int64 `json:"userId"`
}

// AudioFromFileMultipartBody defines parameters for AudioFromFile.
type AudioFromFileMultipartBody struct {
	FileLanguageId string             `json:"fileLanguageId"`
	FilePath       openapi_types.File `json:"filePath"`
	FromLanguageId string             `json:"fromLanguageId"`
	TitleName      string             `json:"titleName"`
	ToLanguageId   string             `json:"toLanguageId"`
}

// AudioFromTitleJSONBody defines parameters for AudioFromTitle.
type AudioFromTitleJSONBody struct {
	FromLanguageId int16 `json:"fromLanguageId"`
	TitleId        int64 `json:"titleId"`
	ToLanguageId   int16 `json:"toLanguageId"`
}

// GetPhrasesParams defines parameters for GetPhrases.
type GetPhrasesParams struct {
	// Limit maximum number of results to return
	Limit *int32 `form:"limit,omitempty" json:"limit,omitempty"`
}

// FindTitlesParams defines parameters for FindTitles.
type FindTitlesParams struct {
	// Similarity find titles similar to
	Similarity string `form:"similarity" json:"similarity"`

	// Limit maximum number of results to return
	Limit int32 `form:"limit" json:"limit"`
}

// AddTitleMultipartBody defines parameters for AddTitle.
type AddTitleMultipartBody struct {
	FilePath   openapi_types.File `json:"filePath"`
	LanguageId string             `json:"languageId"`
	TitleName  string             `json:"titleName"`
}

// AudioFromFileMultipartRequestBody defines body for AudioFromFile for multipart/form-data ContentType.
type AudioFromFileMultipartRequestBody AudioFromFileMultipartBody

// AudioFromTitleJSONRequestBody defines body for AudioFromTitle for application/json ContentType.
type AudioFromTitleJSONRequestBody AudioFromTitleJSONBody

// AddTitleMultipartRequestBody defines body for AddTitle for multipart/form-data ContentType.
type AddTitleMultipartRequestBody AddTitleMultipartBody

// TitlesTranslateJSONRequestBody defines body for TitlesTranslate for application/json ContentType.
type TitlesTranslateJSONRequestBody = TitlesTranslateRequest

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = NewUser

// LoginUserJSONRequestBody defines body for LoginUser for application/json ContentType.
type LoginUserJSONRequestBody = UserLogin

// UpdateUserApplicationJSONPatchPlusJSONRequestBody defines body for UpdateUser for application/json-patch+json ContentType.
type UpdateUserApplicationJSONPatchPlusJSONRequestBody = PatchRequest

// AddUserPermissionJSONRequestBody defines body for AddUserPermission for application/json ContentType.
type AddUserPermissionJSONRequestBody = NewUserPermission

// UpdateUsersPhrasesApplicationJSONPatchPlusJSONRequestBody defines body for UpdateUsersPhrases for application/json-patch+json ContentType.
type UpdateUsersPhrasesApplicationJSONPatchPlusJSONRequestBody = PatchRequest

// AsPatchRequestAddReplaceTest returns the union data inside the PatchRequest_Item as a PatchRequestAddReplaceTest
func (t PatchRequest_Item) AsPatchRequestAddReplaceTest() (PatchRequestAddReplaceTest, error) {
	var body PatchRequestAddReplaceTest
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromPatchRequestAddReplaceTest overwrites any union data inside the PatchRequest_Item as the provided PatchRequestAddReplaceTest
func (t *PatchRequest_Item) FromPatchRequestAddReplaceTest(v PatchRequestAddReplaceTest) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergePatchRequestAddReplaceTest performs a merge with any union data inside the PatchRequest_Item, using the provided PatchRequestAddReplaceTest
func (t *PatchRequest_Item) MergePatchRequestAddReplaceTest(v PatchRequestAddReplaceTest) error {
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

	// (POST /audio/fromfile)
	AudioFromFile(ctx echo.Context) error

	// (POST /audio/fromtitle)
	AudioFromTitle(ctx echo.Context) error
	// Returns list of all available languages
	// (GET /languages)
	GetLanguages(ctx echo.Context) error
	// Returns phrases by title_id
	// (GET /phrases)
	GetPhrases(ctx echo.Context, params GetPhrasesParams) error
	// Returns all titles
	// (GET /titles)
	FindTitles(ctx echo.Context, params FindTitlesParams) error
	// Creates a new title
	// (POST /titles)
	AddTitle(ctx echo.Context) error

	// (POST /titles/translate)
	TitlesTranslate(ctx echo.Context) error
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
	// patches usersphrases resource
	// (PATCH /usersphrases/{phraseId}/{languageId})
	UpdateUsersPhrases(ctx echo.Context, phraseId int64, languageId int16) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// AudioFromFile converts echo context to params.
func (w *ServerInterfaceWrapper) AudioFromFile(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"titles:w"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AudioFromFile(ctx)
	return err
}

// AudioFromTitle converts echo context to params.
func (w *ServerInterfaceWrapper) AudioFromTitle(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"titles:w"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AudioFromTitle(ctx)
	return err
}

// GetLanguages converts echo context to params.
func (w *ServerInterfaceWrapper) GetLanguages(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetLanguages(ctx)
	return err
}

// GetPhrases converts echo context to params.
func (w *ServerInterfaceWrapper) GetPhrases(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPhrasesParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPhrases(ctx, params)
	return err
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

// TitlesTranslate converts echo context to params.
func (w *ServerInterfaceWrapper) TitlesTranslate(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"titles:w"})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.TitlesTranslate(ctx)
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

// UpdateUsersPhrases converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateUsersPhrases(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "phraseId" -------------
	var phraseId int64

	err = runtime.BindStyledParameterWithOptions("simple", "phraseId", ctx.Param("phraseId"), &phraseId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter phraseId: %s", err))
	}

	// ------------- Path parameter "languageId" -------------
	var languageId int16

	err = runtime.BindStyledParameterWithOptions("simple", "languageId", ctx.Param("languageId"), &languageId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter languageId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateUsersPhrases(ctx, phraseId, languageId)
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

	router.POST(baseURL+"/audio/fromfile", wrapper.AudioFromFile)
	router.POST(baseURL+"/audio/fromtitle", wrapper.AudioFromTitle)
	router.GET(baseURL+"/languages", wrapper.GetLanguages)
	router.GET(baseURL+"/phrases", wrapper.GetPhrases)
	router.GET(baseURL+"/titles", wrapper.FindTitles)
	router.POST(baseURL+"/titles", wrapper.AddTitle)
	router.POST(baseURL+"/titles/translate", wrapper.TitlesTranslate)
	router.DELETE(baseURL+"/titles/:id", wrapper.DeleteTitle)
	router.GET(baseURL+"/titles/:id", wrapper.FindTitleByID)
	router.POST(baseURL+"/users", wrapper.CreateUser)
	router.POST(baseURL+"/users/login", wrapper.LoginUser)
	router.DELETE(baseURL+"/users/:id", wrapper.DeleteUser)
	router.GET(baseURL+"/users/:id", wrapper.FindUserByID)
	router.PATCH(baseURL+"/users/:id", wrapper.UpdateUser)
	router.POST(baseURL+"/userspermissions", wrapper.AddUserPermission)
	router.PATCH(baseURL+"/usersphrases/:phraseId/:languageId", wrapper.UpdateUsersPhrases)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"///M1LcwGEEAAA==",
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
