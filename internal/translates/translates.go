package translates

import (
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"cloud.google.com/go/translate"
	"context"
	"errors"
	"fmt"
	"github.com/googleapis/gax-go/v2"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
	"os"
	"strconv"
	"strings"
	"sync"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/util"
	"time"
	"unicode"
)

// TranslateClientX creates an interface for google translate.Translate so it can
// be mocked for testing
type TranslateClientX interface {
	Translate(context.Context, []string, language.Tag, *translate.Options) ([]translate.Translation, error)
}

// TTSClientX creates an interface for google texttospeechpb.SynthesizeSpeech so
// it can be mocked for testing
type TTSClientX interface {
	SynthesizeSpeech(context.Context, *texttospeechpb.SynthesizeSpeechRequest, ...gax.CallOption) (*texttospeechpb.SynthesizeSpeechResponse, error)
}

// TranslateX creates an interface for Translate methods that deal with creating/inserting
// translates and phrases
type TranslateX interface {
	InsertNewPhrases(echo.Context, db.Title, db.Querier, []string) ([]db.Translate, error)
	InsertTranslates(echo.Context, db.Querier, int16, []util.TranslatesReturn) ([]db.Translate, error)
	CreateTTS(echo.Context, db.Querier, db.Title, int16, string) error
	TranslatePhrases(echo.Context, []db.Translate, db.Language) ([]util.TranslatesReturn, error)
}

type Translate struct {
	translateClient TranslateClientX
	ttsClient       TTSClientX
}

func New(trc TranslateClientX, ttsc TTSClientX) *Translate {
	return &Translate{
		translateClient: trc,
		ttsClient:       ttsc,
	}
}

// TranslatePhrases takes a slice of db.Translate{} and a db.Language and returns a slice
// of util.TranslatesReturn to be inserted into the db
func (t *Translate) TranslatePhrases(e echo.Context, ts []db.Translate, dbLang db.Language) ([]util.TranslatesReturn, error) {

	// get language tag to translate to
	langTag, err := language.Parse(dbLang.Tag)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	// concurrently get all the responses from Google Translate
	var wg sync.WaitGroup
	responses := make([]util.TranslatesReturn, len(ts)) // create string slice to hold all the responses
	// create context with cancel, so you can cancel all other requests after any error
	newCtx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure it's called to release resources even if no errors

	for i, nextTranslate := range ts {
		// added intermittent sleep to fix TLS handshake errors on the client side
		if i%50 == 0 && i != 0 {
			time.Sleep(2 * time.Second)
		}
		wg.Add(1)
		//get responses concurrently with go routines
		go t.GetTranslate(e, newCtx, cancel, nextTranslate, &wg, langTag, responses, i)
	}
	wg.Wait()

	if newCtx.Err() != nil {
		e.Logger().Error(newCtx.Err())
		return nil, newCtx.Err()
	}

	return responses, nil
}

// GetTranslate is a helper function for TranslatePhrases that allows concurrent calls to
// google translate.Translate.
// It receives a context.CancelFunc that is invoked on an error so all subsequent calls to
// google translate.Translate can be aborted
func (t *Translate) GetTranslate(e echo.Context,
	ctx context.Context,
	cancel context.CancelFunc,
	phrase db.Translate,
	wg *sync.WaitGroup,
	lang language.Tag,
	responses []util.TranslatesReturn,
	i int,
) {

	defer wg.Done()
	select {
	case <-ctx.Done():
		return // Error somewhere, terminate
	default: // Default to avoid blocking

		resp, err := t.translateClient.Translate(ctx, []string{phrase.Phrase}, lang, nil)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				return
			default:
				e.Logger().Error(fmt.Errorf("error translating text: %s", err))
				cancel()
			}
			return
		}

		if len(resp) == 0 {
			e.Logger().Error(fmt.Errorf("translate returned empty response to text: %s", err))
			cancel()
		}

		//app.Logger.PrintInfo(fmt.Sprintf("response is: %s", resp[0].Text), nil)
		responses[i] = util.TranslatesReturn{
			PhraseId: phrase.PhraseID,
			Text:     resp[0].Text,
		}
	}
}

// CreateTTS is called from api.createAudioFromTitle.
// It checks if the mp3 audio files exist and if not it creates them.
func (t *Translate) CreateTTS(e echo.Context, q db.Querier, title db.Title, voiceId int16, basePath string) error {

	voice, err := q.SelectVoiceById(e.Request().Context(), voiceId)
	if err != nil {
		e.Logger().Error(err)
		return err
	}

	lang, err := q.SelectLanguagesById(e.Request().Context(), voice.LanguageID)
	if err != nil {
		e.Logger().Error(err)
		return err
	}

	// if the audio files already exist no need to request them again
	skip, err := util.PathExists(basePath)
	if err != nil {
		e.Logger().Error(err)
		return err
	}

	// if they do not exist, then request them
	if !skip {
		fromTranslates, err := t.GetOrCreateTranslates(e, q, title, lang)
		if err != nil {
			return err
		}

		err = os.MkdirAll(basePath, 0777)
		if err != nil {
			e.Logger().Error(err)
			return err
		}

		if err = t.TextToSpeech(e, fromTranslates, voice, basePath); err != nil {
			e.Logger().Error(err)
			return err
		}
	}

	return nil
}

// TextToSpeech takes a slice of db.Translate and get the speech mp3's adding them
// to the machines local file system
func (t *Translate) TextToSpeech(e echo.Context, ts []db.Translate, voice db.Voice, bp string) error {

	// set the texttospeec params from the db voice sent in the request
	voiceSelectionParams := &texttospeechpb.VoiceSelectionParams{
		LanguageCode: voice.LanguageCodes[0],
		SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
		Name:         voice.Name,
	}
	if voice.SsmlGender == "FEMALE" {
		voiceSelectionParams.SsmlGender = texttospeechpb.SsmlVoiceGender_FEMALE
	}
	// concurrently get all the audio content from Google text-to-speech
	var wg sync.WaitGroup
	// create context with cancel, so you can cancel all other requests after any error
	newCtx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure it's called to release resources even if no errors

	for i, nextText := range ts {
		// added intermittent sleep to fix TLS handshake errors on the client side
		if i%50 == 0 && i != 0 {
			time.Sleep(2 * time.Second)
		}
		wg.Add(1)
		//get responses concurrently with go routines
		go t.GetSpeech(e, newCtx, cancel, nextText, &wg, voiceSelectionParams, bp)
	}
	wg.Wait()

	if newCtx.Err() != nil {
		e.Logger().Error(newCtx.Err())
		return newCtx.Err()
	}

	return nil
}

// GetSpeech is a helper function for TextToSpeech that is run concurrently.
// it is passed a cancel context, so if one routine fails, the following routines can
// be canceled
func (t *Translate) GetSpeech(
	e echo.Context,
	ctx context.Context,
	cancel context.CancelFunc,
	translate db.Translate,
	wg *sync.WaitGroup,
	params *texttospeechpb.VoiceSelectionParams,
	basePath string) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return // Error somewhere, terminate
	default:
		// Perform the text-to-speech request on the text input with the selected
		// voice parameters and audio file type.
		req := texttospeechpb.SynthesizeSpeechRequest{
			// Set the text input to be synthesized.
			Input: &texttospeechpb.SynthesisInput{
				InputSource: &texttospeechpb.SynthesisInput_Text{Text: translate.Phrase},
			},
			// Build the voice request, select the language code ("en-US") and the SSML
			// voice gender ("neutral").
			Voice: params,
			// Select the type of audio file you want returned.
			AudioConfig: &texttospeechpb.AudioConfig{
				AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			},
		}

		resp, err := t.ttsClient.SynthesizeSpeech(ctx, &req)
		if err != nil {
			e.Logger().Error(fmt.Errorf("error creating Synthesize Speech client: %s", err))
			cancel()
			return
		}

		// The resp AudioContent is binary.
		filename := basePath + strconv.FormatInt(translate.PhraseID, 10)
		err = os.WriteFile(filename, resp.AudioContent, 0644)
		if err != nil {
			e.Logger().Error(fmt.Errorf("error creating translate client: %s", err))
			cancel()
			return
		}
	}
}

// GetOrCreateTranslates checks if the translates for the title already exists in the db.
// If they do not exist, then it creates and returns them.
func (t *Translate) GetOrCreateTranslates(e echo.Context, q db.Querier, title db.Title, lang db.Language) ([]db.Translate, error) {
	// see if translates exist for title for language
	exists, err := q.SelectExistsTranslates(
		e.Request().Context(),
		db.SelectExistsTranslatesParams{
			LanguageID: lang.ID,
			ID:         title.ID,
		})

	// if exists get translates for language
	if exists {
		params := db.SelectTranslatesByTitleIdLangIdParams{
			LanguageID: lang.ID,
			ID:         title.ID,
		}
		translates, err := q.SelectTranslatesByTitleIdLangId(e.Request().Context(), params)
		if err != nil {
			e.Logger().Error(err)
			return nil, err
		}
		return translates, nil
	}

	// if not exists get translates for title original language
	fromTranslates, err := q.SelectTranslatesByTitleIdLangId(
		e.Request().Context(),
		db.SelectTranslatesByTitleIdLangIdParams{
			LanguageID: title.OgLanguageID,
			ID:         title.ID,
		})
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	// create translates for title and to language and return
	translatesReturn, err := t.TranslatePhrases(e, fromTranslates, lang)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}

	dbTranslates, err := t.InsertTranslates(e, q, lang.ID, translatesReturn)
	if err != nil {
		e.Logger().Error(err)
		return nil, err
	}
	return dbTranslates, nil
}

// InsertNewPhrases accepts a slice of phrases and inserts them into the db as translates
func (t *Translate) InsertNewPhrases(e echo.Context, title db.Title, q db.Querier, stringsSlice []string) ([]db.Translate, error) {
	dbTranslates := make([]db.Translate, len(stringsSlice))
	for i, str := range stringsSlice {

		phrase, err := q.InsertPhrases(e.Request().Context(), title.ID)
		if err != nil {
			return nil, err
		}

		insertTranslate, err := q.InsertTranslates(
			e.Request().Context(),
			db.InsertTranslatesParams{
				PhraseID:   phrase.ID,
				LanguageID: title.OgLanguageID,
				Phrase:     str,
				PhraseHint: makeHintString(str),
			})
		if err != nil {
			return nil, err
		}

		dbTranslates[i] = insertTranslate
	}

	return dbTranslates, nil
}

// InsertTranslates accepts a slice of util.TranslatedReturns and inserts them into the db
func (t *Translate) InsertTranslates(e echo.Context, q db.Querier, langId int16, trr []util.TranslatesReturn) ([]db.Translate, error) {
	dbTranslates := make([]db.Translate, len(trr))
	for i, row := range trr {

		// apostrophe's are replaced with &#39; in the response from google translate
		replacedText := strings.ReplaceAll(row.Text, "&#39;", "'")
		insertTranslate, err := q.InsertTranslates(
			e.Request().Context(),
			db.InsertTranslatesParams{
				PhraseID:   row.PhraseId,
				LanguageID: langId,
				Phrase:     replacedText,
				PhraseHint: makeHintString(replacedText),
			})
		if err != nil {
			return nil, err
		}

		dbTranslates[i] = insertTranslate
	}

	return dbTranslates, nil
}

// makeHintString creates a hint string that is the first character of each word of a phrase
// and an underscore for every following character giving the user help when requested.
func makeHintString(s string) string {
	hintString := ""
	words := strings.Fields(s)
	for _, word := range words {
		punctuation := false
		hintString += string(word[0])
		if unicode.IsPunct(rune(word[0])) {
			punctuation = true
		}
		for i := 1; i < len(word); i++ {
			if punctuation {
				hintString += string(word[i])
				punctuation = false
			} else if unicode.IsLetter(rune(word[i])) {
				hintString += "_"
			} else {
				hintString += string(word[i])
			}
		}
		hintString += " "
	}
	return hintString
}
