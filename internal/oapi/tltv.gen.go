// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
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

// AudioFromTitle defines model for AudioFromTitle.
type AudioFromTitle struct {
	FromLanguageId int16  `json:"fromLanguageId"`
	FromVoiceId    *int16 `json:"fromVoiceId,omitempty"`
	TitleId        int64  `json:"titleId"`
	ToLanguageId   int16  `json:"toLanguageId"`
	ToVoiceId      *int16 `json:"toVoiceId,omitempty"`
}

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

// Voice defines model for Voice.
type Voice struct {
	// Id id of voice
	Id int16 `json:"id"`

	// LanguageId id of language
	LanguageId int16 `json:"languageId"`

	// Name the name of the voice
	Name string `json:"name"`

	// NaturalSampleRateHertz the natural sample rate of the voice in hertz
	NaturalSampleRateHertz int16 `json:"naturalSampleRateHertz"`

	// SsmlGender gender of voice MALE|FEMALE
	SsmlGender string `json:"ssmlGender"`
}

// AudioFromFileMultipartBody defines parameters for AudioFromFile.
type AudioFromFileMultipartBody struct {
	FileLanguageId string             `json:"fileLanguageId"`
	FilePath       openapi_types.File `json:"filePath"`
	FromLanguageId string             `json:"fromLanguageId"`
	FromVoiceId    string             `json:"fromVoiceId"`
	TitleName      string             `json:"titleName"`
	ToLanguageId   string             `json:"toLanguageId"`
	ToVoiceId      string             `json:"toVoiceId"`
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

// GetVoicesParams defines parameters for GetVoices.
type GetVoicesParams struct {
	// LanguageId filter by language_id
	LanguageId *int16 `form:"language_id,omitempty" json:"language_id,omitempty"`
}

// AudioFromFileMultipartRequestBody defines body for AudioFromFile for multipart/form-data ContentType.
type AudioFromFileMultipartRequestBody AudioFromFileMultipartBody

// AudioFromTitleJSONRequestBody defines body for AudioFromTitle for application/json ContentType.
type AudioFromTitleJSONRequestBody = AudioFromTitle

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
	// Returns list of all available voices
	// (GET /voices)
	GetVoices(ctx echo.Context, params GetVoicesParams) error
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

	ctx.Set(BearerAuthScopes, []string{})

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

// GetVoices converts echo context to params.
func (w *ServerInterfaceWrapper) GetVoices(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetVoicesParams
	// ------------- Optional query parameter "language_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "language_id", ctx.QueryParams(), &params.LanguageId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter language_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetVoices(ctx, params)
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
	router.GET(baseURL+"/voices", wrapper.GetVoices)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"pDBbbQ4B6PT/FQAA///wrJaFVkYAAA==",
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
