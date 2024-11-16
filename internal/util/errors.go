package util

import "errors"

var (
	ErrVoiceLangIdNoMatch = errors.New("voice id does not match chosen language id")
	ErrOneFile            = errors.New("no need to zip one file")
)
