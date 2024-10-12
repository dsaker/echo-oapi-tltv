package util

// created struct here to not cause cycle between api and internal/mock
type TranslatesReturn struct {
	PhraseId int64
	Text     string
}
