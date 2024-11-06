package util

// avoiding cycle between mock/translates and translates/translates
type TranslatesReturn struct {
	PhraseId int64
	Text     string
}
