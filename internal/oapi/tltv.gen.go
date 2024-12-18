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
	FromVoiceId int16 `json:"fromVoiceId"`

	// Pattern pattern is the pattern used to construct the audio files. You have 3 choices:
	// 1 is standard and should be used if you are at a beginner or intermediate level of language learning
	// 2 is intermediate
	// 3 is advanced and repeats phrases less often and should only be used if you are at an advanced level
	// 4 is review and only repeats each phrase one time and can be used to review already learned phrases
	Pattern *int `json:"pattern,omitempty"`

	// Pause the pause in seconds between phrases in the audio file (default is 4)
	Pause     *int  `json:"pause,omitempty"`
	TitleId   int64 `json:"titleId"`
	ToVoiceId int16 `json:"toVoiceId"`
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
	// FileLanguageId the original language of the file you are uploading
	FileLanguageId string             `json:"fileLanguageId"`
	FilePath       openapi_types.File `json:"filePath"`

	// FromVoiceId the language you know
	FromVoiceId string `json:"fromVoiceId"`

	// Pattern pattern is the pattern used to construct the audio files. You have 3 choices:
	// 1 is standard and should be used if you are at a beginner or intermediate level of language learning
	// 2 is intermediate
	// 3 is advanced and repeats phrases less often and should only be used if you are at an advanced level
	// 4 is review and only repeats each phrase one time and can be used to review already learned phrases
	Pattern *string `json:"pattern,omitempty"`

	// Pause the pause in seconds between phrases in the audiofile (default is 4)
	Pause *string `json:"pause,omitempty"`

	// TitleName choose a descriptive title that includes to and from languages
	TitleName string `json:"titleName"`

	// ToVoiceId the language you want to learn
	ToVoiceId string `json:"toVoiceId"`
}

// GetLanguagesParams defines parameters for GetLanguages.
type GetLanguagesParams struct {
	// Similarity find titles similar to
	Similarity *string `form:"similarity,omitempty" json:"similarity,omitempty"`
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
	// LanguageId filter by languageId
	LanguageId *int16 `form:"languageId,omitempty" json:"languageId,omitempty"`
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
	GetLanguages(ctx echo.Context, params GetLanguagesParams) error
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

	// Parameter object where we will unmarshal all parameters from the context
	var params GetLanguagesParams
	// ------------- Optional query parameter "similarity" -------------

	err = runtime.BindQueryParameter("form", true, false, "similarity", ctx.QueryParams(), &params.Similarity)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter similarity: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetLanguages(ctx, params)
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
	// ------------- Optional query parameter "languageId" -------------

	err = runtime.BindQueryParameter("form", true, false, "languageId", ctx.QueryParams(), &params.LanguageId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter languageId: %s", err))
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
