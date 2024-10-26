package clients

import (
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"cloud.google.com/go/translate"
	"context"
	"github.com/googleapis/gax-go/v2"
	"golang.org/x/text/language"
)

type TranslateClientX interface {
	Translate(context.Context, []string, language.Tag, *translate.Options) ([]translate.Translation, error)
}

type TTSClientX interface {
	SynthesizeSpeech(context.Context, *texttospeechpb.SynthesizeSpeechRequest, ...gax.CallOption) (*texttospeechpb.SynthesizeSpeechResponse, error)
}
