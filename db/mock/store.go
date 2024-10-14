// Code generated by MockGen. DO NOT EDIT.
// Source: talkliketv.click/tltv/db/sqlc (interfaces: Querier)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	db "talkliketv.click/tltv/db/sqlc"
)

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// DeleteTitleById mocks base method.
func (m *MockQuerier) DeleteTitleById(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTitleById", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTitleById indicates an expected call of DeleteTitleById.
func (mr *MockQuerierMockRecorder) DeleteTitleById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTitleById", reflect.TypeOf((*MockQuerier)(nil).DeleteTitleById), arg0, arg1)
}

// DeleteUserById mocks base method.
func (m *MockQuerier) DeleteUserById(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserById", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserById indicates an expected call of DeleteUserById.
func (mr *MockQuerierMockRecorder) DeleteUserById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserById", reflect.TypeOf((*MockQuerier)(nil).DeleteUserById), arg0, arg1)
}

// DeleteUserPermissionById mocks base method.
func (m *MockQuerier) DeleteUserPermissionById(arg0 context.Context, arg1 db.DeleteUserPermissionByIdParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserPermissionById", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserPermissionById indicates an expected call of DeleteUserPermissionById.
func (mr *MockQuerierMockRecorder) DeleteUserPermissionById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserPermissionById", reflect.TypeOf((*MockQuerier)(nil).DeleteUserPermissionById), arg0, arg1)
}

// InsertPhrases mocks base method.
func (m *MockQuerier) InsertPhrases(arg0 context.Context, arg1 int64) (db.Phrase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertPhrases", arg0, arg1)
	ret0, _ := ret[0].(db.Phrase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertPhrases indicates an expected call of InsertPhrases.
func (mr *MockQuerierMockRecorder) InsertPhrases(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertPhrases", reflect.TypeOf((*MockQuerier)(nil).InsertPhrases), arg0, arg1)
}

// InsertTitle mocks base method.
func (m *MockQuerier) InsertTitle(arg0 context.Context, arg1 db.InsertTitleParams) (db.Title, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTitle", arg0, arg1)
	ret0, _ := ret[0].(db.Title)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertTitle indicates an expected call of InsertTitle.
func (mr *MockQuerierMockRecorder) InsertTitle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTitle", reflect.TypeOf((*MockQuerier)(nil).InsertTitle), arg0, arg1)
}

// InsertTranslates mocks base method.
func (m *MockQuerier) InsertTranslates(arg0 context.Context, arg1 db.InsertTranslatesParams) (db.Translate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTranslates", arg0, arg1)
	ret0, _ := ret[0].(db.Translate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertTranslates indicates an expected call of InsertTranslates.
func (mr *MockQuerierMockRecorder) InsertTranslates(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTranslates", reflect.TypeOf((*MockQuerier)(nil).InsertTranslates), arg0, arg1)
}

// InsertUser mocks base method.
func (m *MockQuerier) InsertUser(arg0 context.Context, arg1 db.InsertUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertUser indicates an expected call of InsertUser.
func (mr *MockQuerierMockRecorder) InsertUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUser", reflect.TypeOf((*MockQuerier)(nil).InsertUser), arg0, arg1)
}

// InsertUserPermission mocks base method.
func (m *MockQuerier) InsertUserPermission(arg0 context.Context, arg1 db.InsertUserPermissionParams) (db.UsersPermission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUserPermission", arg0, arg1)
	ret0, _ := ret[0].(db.UsersPermission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertUserPermission indicates an expected call of InsertUserPermission.
func (mr *MockQuerierMockRecorder) InsertUserPermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUserPermission", reflect.TypeOf((*MockQuerier)(nil).InsertUserPermission), arg0, arg1)
}

// ListLanguages mocks base method.
func (m *MockQuerier) ListLanguages(arg0 context.Context) ([]db.Language, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListLanguages", arg0)
	ret0, _ := ret[0].([]db.Language)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListLanguages indicates an expected call of ListLanguages.
func (mr *MockQuerierMockRecorder) ListLanguages(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListLanguages", reflect.TypeOf((*MockQuerier)(nil).ListLanguages), arg0)
}

// ListTitles mocks base method.
func (m *MockQuerier) ListTitles(arg0 context.Context, arg1 db.ListTitlesParams) ([]db.ListTitlesRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTitles", arg0, arg1)
	ret0, _ := ret[0].([]db.ListTitlesRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTitles indicates an expected call of ListTitles.
func (mr *MockQuerierMockRecorder) ListTitles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTitles", reflect.TypeOf((*MockQuerier)(nil).ListTitles), arg0, arg1)
}

// ListTitlesByOgLanguage mocks base method.
func (m *MockQuerier) ListTitlesByOgLanguage(arg0 context.Context, arg1 db.ListTitlesByOgLanguageParams) ([]db.ListTitlesByOgLanguageRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTitlesByOgLanguage", arg0, arg1)
	ret0, _ := ret[0].([]db.ListTitlesByOgLanguageRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTitlesByOgLanguage indicates an expected call of ListTitlesByOgLanguage.
func (mr *MockQuerierMockRecorder) ListTitlesByOgLanguage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTitlesByOgLanguage", reflect.TypeOf((*MockQuerier)(nil).ListTitlesByOgLanguage), arg0, arg1)
}

// SelectExistsTranslates mocks base method.
func (m *MockQuerier) SelectExistsTranslates(arg0 context.Context, arg1 db.SelectExistsTranslatesParams) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectExistsTranslates", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectExistsTranslates indicates an expected call of SelectExistsTranslates.
func (mr *MockQuerierMockRecorder) SelectExistsTranslates(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectExistsTranslates", reflect.TypeOf((*MockQuerier)(nil).SelectExistsTranslates), arg0, arg1)
}

// SelectLanguagesById mocks base method.
func (m *MockQuerier) SelectLanguagesById(arg0 context.Context, arg1 int16) (db.Language, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectLanguagesById", arg0, arg1)
	ret0, _ := ret[0].(db.Language)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectLanguagesById indicates an expected call of SelectLanguagesById.
func (mr *MockQuerierMockRecorder) SelectLanguagesById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectLanguagesById", reflect.TypeOf((*MockQuerier)(nil).SelectLanguagesById), arg0, arg1)
}

// SelectPermissionByCode mocks base method.
func (m *MockQuerier) SelectPermissionByCode(arg0 context.Context, arg1 string) (db.Permission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectPermissionByCode", arg0, arg1)
	ret0, _ := ret[0].(db.Permission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectPermissionByCode indicates an expected call of SelectPermissionByCode.
func (mr *MockQuerierMockRecorder) SelectPermissionByCode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectPermissionByCode", reflect.TypeOf((*MockQuerier)(nil).SelectPermissionByCode), arg0, arg1)
}

// SelectPhrasesFromTranslates mocks base method.
func (m *MockQuerier) SelectPhrasesFromTranslates(arg0 context.Context, arg1 db.SelectPhrasesFromTranslatesParams) ([]db.SelectPhrasesFromTranslatesRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectPhrasesFromTranslates", arg0, arg1)
	ret0, _ := ret[0].([]db.SelectPhrasesFromTranslatesRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectPhrasesFromTranslates indicates an expected call of SelectPhrasesFromTranslates.
func (mr *MockQuerierMockRecorder) SelectPhrasesFromTranslates(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectPhrasesFromTranslates", reflect.TypeOf((*MockQuerier)(nil).SelectPhrasesFromTranslates), arg0, arg1)
}

// SelectPhrasesFromTranslatesWithCorrect mocks base method.
func (m *MockQuerier) SelectPhrasesFromTranslatesWithCorrect(arg0 context.Context, arg1 db.SelectPhrasesFromTranslatesWithCorrectParams) ([]db.SelectPhrasesFromTranslatesWithCorrectRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectPhrasesFromTranslatesWithCorrect", arg0, arg1)
	ret0, _ := ret[0].([]db.SelectPhrasesFromTranslatesWithCorrectRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectPhrasesFromTranslatesWithCorrect indicates an expected call of SelectPhrasesFromTranslatesWithCorrect.
func (mr *MockQuerierMockRecorder) SelectPhrasesFromTranslatesWithCorrect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectPhrasesFromTranslatesWithCorrect", reflect.TypeOf((*MockQuerier)(nil).SelectPhrasesFromTranslatesWithCorrect), arg0, arg1)
}

// SelectTitleById mocks base method.
func (m *MockQuerier) SelectTitleById(arg0 context.Context, arg1 int64) (db.Title, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectTitleById", arg0, arg1)
	ret0, _ := ret[0].(db.Title)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectTitleById indicates an expected call of SelectTitleById.
func (mr *MockQuerierMockRecorder) SelectTitleById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectTitleById", reflect.TypeOf((*MockQuerier)(nil).SelectTitleById), arg0, arg1)
}

// SelectTranslatesByTitleIdLangId mocks base method.
func (m *MockQuerier) SelectTranslatesByTitleIdLangId(arg0 context.Context, arg1 db.SelectTranslatesByTitleIdLangIdParams) ([]db.SelectTranslatesByTitleIdLangIdRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectTranslatesByTitleIdLangId", arg0, arg1)
	ret0, _ := ret[0].([]db.SelectTranslatesByTitleIdLangIdRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectTranslatesByTitleIdLangId indicates an expected call of SelectTranslatesByTitleIdLangId.
func (mr *MockQuerierMockRecorder) SelectTranslatesByTitleIdLangId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectTranslatesByTitleIdLangId", reflect.TypeOf((*MockQuerier)(nil).SelectTranslatesByTitleIdLangId), arg0, arg1)
}

// SelectUserById mocks base method.
func (m *MockQuerier) SelectUserById(arg0 context.Context, arg1 int64) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectUserById", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectUserById indicates an expected call of SelectUserById.
func (mr *MockQuerierMockRecorder) SelectUserById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectUserById", reflect.TypeOf((*MockQuerier)(nil).SelectUserById), arg0, arg1)
}

// SelectUserByName mocks base method.
func (m *MockQuerier) SelectUserByName(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectUserByName", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectUserByName indicates an expected call of SelectUserByName.
func (mr *MockQuerierMockRecorder) SelectUserByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectUserByName", reflect.TypeOf((*MockQuerier)(nil).SelectUserByName), arg0, arg1)
}

// SelectUserPermissions mocks base method.
func (m *MockQuerier) SelectUserPermissions(arg0 context.Context, arg1 int64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectUserPermissions", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectUserPermissions indicates an expected call of SelectUserPermissions.
func (mr *MockQuerierMockRecorder) SelectUserPermissions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectUserPermissions", reflect.TypeOf((*MockQuerier)(nil).SelectUserPermissions), arg0, arg1)
}

// SelectUsersPhrasesByCorrect mocks base method.
func (m *MockQuerier) SelectUsersPhrasesByCorrect(arg0 context.Context, arg1 db.SelectUsersPhrasesByCorrectParams) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectUsersPhrasesByCorrect", arg0, arg1)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectUsersPhrasesByCorrect indicates an expected call of SelectUsersPhrasesByCorrect.
func (mr *MockQuerierMockRecorder) SelectUsersPhrasesByCorrect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectUsersPhrasesByCorrect", reflect.TypeOf((*MockQuerier)(nil).SelectUsersPhrasesByCorrect), arg0, arg1)
}

// SelectUsersPhrasesByIds mocks base method.
func (m *MockQuerier) SelectUsersPhrasesByIds(arg0 context.Context, arg1 db.SelectUsersPhrasesByIdsParams) (db.UsersPhrase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectUsersPhrasesByIds", arg0, arg1)
	ret0, _ := ret[0].(db.UsersPhrase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectUsersPhrasesByIds indicates an expected call of SelectUsersPhrasesByIds.
func (mr *MockQuerierMockRecorder) SelectUsersPhrasesByIds(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectUsersPhrasesByIds", reflect.TypeOf((*MockQuerier)(nil).SelectUsersPhrasesByIds), arg0, arg1)
}

// UpdateUserById mocks base method.
func (m *MockQuerier) UpdateUserById(arg0 context.Context, arg1 db.UpdateUserByIdParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserById", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserById indicates an expected call of UpdateUserById.
func (mr *MockQuerierMockRecorder) UpdateUserById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserById", reflect.TypeOf((*MockQuerier)(nil).UpdateUserById), arg0, arg1)
}

// UpdateUsersPhrasesByThreeIds mocks base method.
func (m *MockQuerier) UpdateUsersPhrasesByThreeIds(arg0 context.Context, arg1 db.UpdateUsersPhrasesByThreeIdsParams) (db.UsersPhrase, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUsersPhrasesByThreeIds", arg0, arg1)
	ret0, _ := ret[0].(db.UsersPhrase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUsersPhrasesByThreeIds indicates an expected call of UpdateUsersPhrasesByThreeIds.
func (mr *MockQuerierMockRecorder) UpdateUsersPhrasesByThreeIds(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUsersPhrasesByThreeIds", reflect.TypeOf((*MockQuerier)(nil).UpdateUsersPhrasesByThreeIds), arg0, arg1)
}
