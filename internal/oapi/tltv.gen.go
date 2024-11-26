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

	"H4sIAAAAAAAC/+xcbXPbNhL+Kzi2NxcnsiS/JE01c5O4cdK646aexOnNnaX2IHIloiEBBgAlq7but98s",
	"wBdQImXKiXPO3X1RLAIEFrsPdh/sQrnyfBEnggPXyhtcecoPIabmz6M0YOKVFPE50xHgk0SKBKRmYNon",
	"UsS/CObDSYBf4ZLGCfbbe7Lf8SZCxlR7A49xvffE63h6kYD9ClOQ3rLjJVRrkBzfDUD5kiWaCe4N8gbC",
	"FNEhkPxrqiAgWhBfcKVl6mvTSlFKMmERqC75u0hJSGdADogfomhqQIZ8D0dSmvKAyoBQHhAVijQKyBjs",
	"oGxCFiIlVAKhmlAyhinjHCQRkqDAMoaAUQ0kghlERExIRPk0pVN8QiVnfDrk+ziL23vID/ARDWaU+2An",
	"lpAA1YokoaQKFIlAKSImGrgrl+DRokk4Xg5opBnyQ5xFwozB3Axi3s4nAuqH2WxEcCCaxWB6+ZQXU2hR",
	"vB9JoMHCLguCXM4h9zpeTC9ZnMbe4LDjxYzbv/fqLZsqWLerNWaqgDBOFPiCB4qMQc8BeKERxlfMSh4E",
	"MKFppHGZhzuuIHt9R5KDOkk0Qtfi00Xkk8NaRGpRh+f9/tMWeF52PAkfUiYh8AYXxcSdyjZxpxgVY4jx",
	"7+BrFOCllEKu7zRfBDXqNJ2JaavKd7Bfu7oYlKLTxoHy5uJVpSXj07WVZRPm3UfLjnea7YZ1yVmwPh0L",
	"3B3ktfIVkTNDdTQr5cqIKyvoeJpO118tNrGmU7sRJkKSqRDTCModHoMORaBu1AtD87oy0KlRzmuYNzlQ",
	"FgGnccNOMdjXgqRJJGjgqmnMOJWLumWKaW6KkxrF523EWsBAtJ36db6A6oCvaQxmpBCK0TYrKe9VEbRT",
	"aiJT2DsFNdsAYsqiGvjiYxQjxbc65cbN/3puXuz6IvacqOM9eDa4oLt/9He//dNXX/95mPb7+0/+8vBR",
	"76/Pfv3tn1fXy3/tjh49eDYYDrs3dtt5eD004w2Hl/29Xfx8ih9j/PDxA/Dh3mQ4vNzfw48D/P4Y2x8H",
	"+Oc3k9H1cDgcOiN8WzPCN5PRzsOht/P8wbNBKf+o/HN39DB/uPNsOOzuPLqhz/VweGEH23980d99PLre",
	"v+jvHo4usPn6or83emb+NB/PdnDIq4Nly+7X6zMOWqvJaIjix8GWGnq0MxyOduo2SP12Q7jxDMsIoi75",
	"KVUawyONkpDyNAbJ/G4FW2n2jo1Hp8CnOsRYYSJS/vWxi7dfUQVHu/9ALTz8ulY6mG/avycVx4m+wYRp",
	"V6zdvVb7ebOfsPNwqtkMXK+67SwJVWouZM0MZ1lL7a5N6OPH8+ey4vKKoTZp+2md4y8JQN0iTXMrTdZS",
	"hhXvluHB+qmOK3RJB1Y8X9Xkjvs7AxkzpYywq44wKdrqVla+mTn6sns7b48WqRsYxcqGzIy2rYKykTvV",
	"FeCqz6j2wzfwIQWlDXHQEJu1Cg4/T7zBxZX3tYSJN/C+6pUnll52XOm5rx8FwRtIIurDOQ62xOEzuaiU",
	"dOGtTLfSf3Dl0SBguGYanTl6n9BIQWfFFCJZ19N5CAT7UPyO4EpAoqaMA+HIVS9wCg91Y+ZFxeHUI3cb",
	"0CDokKxDB08ipksNxBOK4F8V4oj8+Pbn1+RMmCMJHqLCqgPr5TBdG3BGoxTql2WacEmudLlw1fE5zM0M",
	"TvRdhYOdqINKzJZRR4rPzMlgG25pzxJt8LnBQdihYjFjcAukM3fXIwILDkijqAWgC9a47LRZ9zvOPqQF",
	"sXPo2C3kHi0LedW5pFxFVIOzNavS3BCz2FrM0vmQRIst2Gfd4NZ7s6onL8fHk9ctVFD66nXnXOhDrWsi",
	"qqihTXwsgL2+p03TD4zrDc0tz7Uryyve7bgSF9JU5sYl51y8NW7NC7eE7S3jSoFanPtUTFld0LwTNrJm",
	"moIX3sgym+jk5hOU07EQo7LyN6ASwesc5u9zXbOHxHvgjkibZ8ch8tlKnuFOuRVMHJLzEYDZkuBswo6y",
	"webG3f3xGQ27y14IKTHMfcIhTz5HWGwdXprZ5KchkqWz3uzXckWjoU0Gbhs+MTMvbJWo+lQwac4PcSf5",
	"kgtYc97VqaTRW+Nm3lANP4DUfzQNaPoSZToTiSHUHZ8wTkLzeivJlYqj74EHNnpUp5ua54VqyU9Hpy+v",
	"X73Ef7bLs1XM/UIEoLzKzJ38TNagiBGOrsBPJdOLt+icLBa+AypBHqWWWY/Nt1f5mn/82znOYnp7g6y1",
	"lDvUOvGWODDjE2Hzt1xTu8uzHJYXpEov5nTB4bkvYp8q3eWgc3kH3jG2k7f0vVXmChOn0ftT9h7OfzEF",
	"hvVqBKFJEjHfHkACUGzKbZo/hCgxe04RMQPpixhsjQU5DU1JVosA7osUzw0QkDnTIRE6BFnOQ5NEdcmJ",
	"JmIywbEoumCFByb2hylMZGLAZQKSAfeBkPGC0CgSc3xuJdCC+KEQyoqgEvDZhPlFJUCHsCBzyjV2nAg/",
	"VUTwLjFu2pQvbGaUEEo0XGqbMDXSOrWEhEo6lTQJCUK2Y2ogWTkkwRUxbg9XMDMlGPL2zbkZqGOKJLYU",
	"4ehyzqKITIGD2R2UKEAdkJ/ODtxClBGY+ixi2laNMnXoUIp0GpKIKQ34pEvIUXHWjBYdQskcxpUZTVGJ",
	"BDCDSCQxcN0h8xAkZDo0AikhTDGHjrNkhpgSfAuXADykqH8dAnMsqN6zKFKFRBJoYHDDg4JA43ek0Nmr",
	"WT6oyuXxpEHlFHTxuGtqRRHzISMEGaCPEuqHQPa7fa/jpTLKNooa9Hrz+bxLTXNXyGkve1f1Tk9evHz9",
	"9uXufrffDXUcOZnocgvMvI43A2kTJd5et9/tmwxXApwmzBt4B+aRPV+avd0zlurh0tBaJhAIVRODfQnW",
	"yjX7qyxPGRVZKOIeu9QIJyUtHrvmgemXFdxcTCrsWoNIo8Iig4ChpCzGvmIm8Ep7HvtOBIvcw4A9LMRp",
	"pFlCpe4h4ncDqmlZ2K0vQmw6wKGJhWRTxmlUKiKLC2ZheXXS6gBdtktsD7+ti0v44lmWt2hR11gpNa9L",
	"WAiGwrznYl6RYe/JfkPi5P8F6C+oAO0YtN6cn6bqXF90dhDdogS9koB+XcvisuBHSfF4BnkeI6SaMO5H",
	"aWDDCWrMeJrc8Kq2zCna75I8sq5lvr39/tN2hbzXllutuJDmkrez79cTfcs1moPDlOOiJRYiLQKRFnya",
	"ooMV1T6bF9olSCDI9y/PSa9Qpan7UhOWK3uLBcqGs3LlWqZgVGFPvMaF7vf7Kz7YCeC935VN4ZcOuMht",
	"bzohO5mmtfz1uqoyJ4QAMfDN5bPc0SB5KxE3SWYvKdQIkXLke76GgEDWpyTXJiXg0uoMQ2owN+fvZccN",
	"zEXB+aMjc1Fxrg+n51lzczy9vaJWJmlAeLY5muHtdLgB2186VreES+kKB1feFGpw8gZ0KrkqdjaNIkJn",
	"lEWGKbuutAqP70GfOo1I2GLQIJURa8WGDDmzkY0oFrOISpvZLh1q9tjDE6E38D6kYGhOxoyzVlx0x9H8",
	"qv8dfQ5LFjd5WtixPA1ij/vkctI4RiLZ3v4GTkmZ+tsIJo6HzYuIxUyPCgIxXlgM/MayE3Mk5qA08W3e",
	"iShfSMNAzLHtt+KMi4LUwS/PQ94AvoyHEJ7GY5tQkaDSSCtLqVBgF4nIVuogaBZTQd+N98k+Dx6z2l8L",
	"NGYU8z6jsAYrFnnWe9wIPERv5mjyczlRC6UhtnyxyIeHVBHq+0jtjSeqousV44Gt7n1e11YNTJtcXeeO",
	"Yd4syT2BfV76vRH19rRwn0Ffgtac0Wrp3AtD5xShhMPc9u6S49SKjF7SQhBPoFxom0aEYJ3UBcHNdO4W",
	"6ZEtEhXV9P/mw2D7A1a0eley9RHqPK+K5/dNPiU3bIHgLwCx7QinC+0auLp+vFfcQGg+wBRdyqhgXHp9",
	"pq24pJY/TKSYsQACJBQZ1MkYsb66JVZucdzRQafhrkiT9U+C6j0NygP3zN18YYRpG37u/RHH4KI08Zd5",
	"Gs/QfMWCpQVwBLomi2Wf425QjE+jPIM1pgoCImxe7eSYqBSXWOO1j837uePeyEdOjivbIZMoi/Xm/lgR",
	"6lnQOs7X15vX4/xh0wUkK0dwn6LvcWGUzBoLcnKMAm4mmKu2K2x6ctxMI79bmNZtLDcB7YefzXD/i1Ft",
	"nYdVkYAb3BxFm2NUNczZe/EFKSOmkG1cd34jaZWfdYfc/DSjY64q2IJqfmXJ9M2t3K2pd9nJ39nbIXcR",
	"s4qLausaNfecW3OmvU8mUpM85kxnU57GyRxaRK/sMD6jEQtIFny7xPw4J679HcN9CkMbmFWqsps/Fqm9",
	"qLjYV4tX4EEiWJYPxa7dHGHKwagDwjXQmftzd4i58nZig5XrRPx8nH39CmETGI167yetqeDJrIbQNSRt",
	"zWnMqrekNBmQWsTFNHM4/3lCYyS513zGWqIlnama7SY2gwZrT2Zyo32pXGZjtLnfTKaKgYSiAWpudWs/",
	"JJQTuGRK57fO1gz/LgnobfaqnfUOzd42/OwaSR5tZ4PKz62Wy+VdhpUNTOveYq0BPGUIKa+cb2DQR0GQ",
	"c5myv2GW7s9MzX2haSTGNCI0iBkniWQzFsEU6lKcK/fm75Qdu/fzG3jyWXVhweenLDW/Q2iCmmOFLyIn",
	"Y1ExMKgo8jIWfzZr2LvKr70ve1dlmnhpf/BS6xdl5kbdYUiWQW52jqplQTIXJ3eVxQybXaZzef9jHGen",
	"6b942FqgSsq9nUj1Pzb57/LlFSg07LJCw/fRsRsFwwr6JSiRSj/L4pufJLS/U2K7d8kpoB83FUcyjih/",
	"b2+VVW4c2L5ESDJhkbYkxsl8d+tuAvxixbmxVLs6nkFubRXU7bA1lO860W5/nNMix25/OfKl3j3JQLbc",
	"HADQ5f87AAD//wgaeIyiSwAA",
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
