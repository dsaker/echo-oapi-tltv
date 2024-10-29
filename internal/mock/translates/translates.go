// Code generated by MockGen. DO NOT EDIT.
// Source: internal/translates/translates.go

// Package mockt is a generated GoMock package.
package mockt

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo/v4"
	db "talkliketv.click/tltv/db/sqlc"
	util "talkliketv.click/tltv/internal/util"
)

// MockTranslateX is a mock of TranslateX interface.
type MockTranslateX struct {
	ctrl     *gomock.Controller
	recorder *MockTranslateXMockRecorder
}

// MockTranslateXMockRecorder is the mock recorder for MockTranslateX.
type MockTranslateXMockRecorder struct {
	mock *MockTranslateX
}

// NewMockTranslateX creates a new mock instance.
func NewMockTranslateX(ctrl *gomock.Controller) *MockTranslateX {
	mock := &MockTranslateX{ctrl: ctrl}
	mock.recorder = &MockTranslateXMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTranslateX) EXPECT() *MockTranslateXMockRecorder {
	return m.recorder
}

// CreateTTS mocks base method.
func (m *MockTranslateX) CreateTTS(arg0 echo.Context, arg1 db.Querier, arg2 db.Language, arg3 db.Title, arg4 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTTS", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTTS indicates an expected call of CreateTTS.
func (mr *MockTranslateXMockRecorder) CreateTTS(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTTS", reflect.TypeOf((*MockTranslateX)(nil).CreateTTS), arg0, arg1, arg2, arg3, arg4)
}

// CreateTTSForLang mocks base method.
func (m *MockTranslateX) CreateTTSForLang(arg0 echo.Context, arg1 db.Querier, arg2 db.Language, arg3 db.Title, arg4 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTTSForLang", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTTSForLang indicates an expected call of CreateTTSForLang.
func (mr *MockTranslateXMockRecorder) CreateTTSForLang(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTTSForLang", reflect.TypeOf((*MockTranslateX)(nil).CreateTTSForLang), arg0, arg1, arg2, arg3, arg4)
}

// InsertNewPhrases mocks base method.
func (m *MockTranslateX) InsertNewPhrases(arg0 echo.Context, arg1 db.Title, arg2 db.Querier, arg3 []string) ([]db.Translate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNewPhrases", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]db.Translate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertNewPhrases indicates an expected call of InsertNewPhrases.
func (mr *MockTranslateXMockRecorder) InsertNewPhrases(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNewPhrases", reflect.TypeOf((*MockTranslateX)(nil).InsertNewPhrases), arg0, arg1, arg2, arg3)
}

// InsertTranslates mocks base method.
func (m *MockTranslateX) InsertTranslates(arg0 echo.Context, arg1 db.Querier, arg2 int16, arg3 []util.TranslatesReturn) ([]db.Translate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTranslates", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]db.Translate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertTranslates indicates an expected call of InsertTranslates.
func (mr *MockTranslateXMockRecorder) InsertTranslates(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTranslates", reflect.TypeOf((*MockTranslateX)(nil).InsertTranslates), arg0, arg1, arg2, arg3)
}

// TranslatePhrases mocks base method.
func (m *MockTranslateX) TranslatePhrases(arg0 echo.Context, arg1 []db.Translate, arg2 db.Language) ([]util.TranslatesReturn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TranslatePhrases", arg0, arg1, arg2)
	ret0, _ := ret[0].([]util.TranslatesReturn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TranslatePhrases indicates an expected call of TranslatePhrases.
func (mr *MockTranslateXMockRecorder) TranslatePhrases(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TranslatePhrases", reflect.TypeOf((*MockTranslateX)(nil).TranslatePhrases), arg0, arg1, arg2)
}
