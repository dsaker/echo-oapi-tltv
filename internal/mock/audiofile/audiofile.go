// Code generated by MockGen. DO NOT EDIT.
// Source: internal/audio/audiofile/audiofile.go

// Package mocka is a generated GoMock package.
package mocka

import (
	multipart "mime/multipart"
	os "os"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo/v4"
	db "talkliketv.click/tltv/db/sqlc"
)

// MockAudioFileX is a mock of AudioFileX interface.
type MockAudioFileX struct {
	ctrl     *gomock.Controller
	recorder *MockAudioFileXMockRecorder
}

// MockAudioFileXMockRecorder is the mock recorder for MockAudioFileX.
type MockAudioFileXMockRecorder struct {
	mock *MockAudioFileX
}

// NewMockAudioFileX creates a new mock instance.
func NewMockAudioFileX(ctrl *gomock.Controller) *MockAudioFileX {
	mock := &MockAudioFileX{ctrl: ctrl}
	mock.recorder = &MockAudioFileXMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAudioFileX) EXPECT() *MockAudioFileXMockRecorder {
	return m.recorder
}

// BuildAudioInputFiles mocks base method.
func (m *MockAudioFileX) BuildAudioInputFiles(arg0 echo.Context, arg1 []int64, arg2 db.Title, arg3, arg4, arg5, arg6 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildAudioInputFiles", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuildAudioInputFiles indicates an expected call of BuildAudioInputFiles.
func (mr *MockAudioFileXMockRecorder) BuildAudioInputFiles(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildAudioInputFiles", reflect.TypeOf((*MockAudioFileX)(nil).BuildAudioInputFiles), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// CreateMp3ZipWithFfmpeg mocks base method.
func (m *MockAudioFileX) CreateMp3ZipWithFfmpeg(arg0 echo.Context, arg1 db.Title, arg2 string) (*os.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMp3ZipWithFfmpeg", arg0, arg1, arg2)
	ret0, _ := ret[0].(*os.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMp3ZipWithFfmpeg indicates an expected call of CreateMp3ZipWithFfmpeg.
func (mr *MockAudioFileXMockRecorder) CreateMp3ZipWithFfmpeg(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMp3ZipWithFfmpeg", reflect.TypeOf((*MockAudioFileX)(nil).CreateMp3ZipWithFfmpeg), arg0, arg1, arg2)
}

// GetLines mocks base method.
func (m *MockAudioFileX) GetLines(arg0 echo.Context, arg1 multipart.File) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLines", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLines indicates an expected call of GetLines.
func (mr *MockAudioFileXMockRecorder) GetLines(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLines", reflect.TypeOf((*MockAudioFileX)(nil).GetLines), arg0, arg1)
}