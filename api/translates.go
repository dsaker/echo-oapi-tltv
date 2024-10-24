package api

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"cloud.google.com/go/translate"
	"context"
	"errors"
	"fmt"
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

type TranslateX interface {
	InsertNewPhrases(echo.Context, db.Title, db.Querier, []string) ([]db.Translate, error)
	TextToSpeech(echo.Context, []db.Translate, string, string) error
	TranslatePhrases(echo.Context, []db.Translate, db.Language) ([]util.TranslatesReturn, error)
	InsertTranslates(echo.Context, db.Querier, int16, []util.TranslatesReturn) ([]db.Translate, error)
	CreateTTS(echo.Context, db.Querier, db.Language, db.Title, string) error
}

type Translate struct {
}

func (t *Translate) InsertNewPhrases(eCtx echo.Context, title db.Title, q db.Querier, stringsSlice []string) ([]db.Translate, error) {
	dbTranslates := make([]db.Translate, len(stringsSlice))
	for i, str := range stringsSlice {

		phrase, err := q.InsertPhrases(eCtx.Request().Context(), title.ID)
		if err != nil {
			return nil, err
		}

		insertTranslate, err := q.InsertTranslates(
			eCtx.Request().Context(),
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

func (t *Translate) InsertTranslates(eCtx echo.Context, q db.Querier, langId int16, trr []util.TranslatesReturn) ([]db.Translate, error) {
	dbTranslates := make([]db.Translate, len(trr))
	for i, row := range trr {

		// apostrophe's are replaced with &#39; in the response from google translate
		replacedText := strings.ReplaceAll(row.Text, "&#39;", "'")
		insertTranslate, err := q.InsertTranslates(
			eCtx.Request().Context(),
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

func (t *Translate) TextToSpeech(eCtx echo.Context, translatesSlice []db.Translate, basepath, tag string) error {

	// concurrently get all the audio content from Google text-to-speech
	var wg sync.WaitGroup
	// create context with cancel, so you can cancel all other requests after any error
	newCtx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure it's called to release resources even if no errors

	for i, nextText := range translatesSlice {
		// added intermittent sleep to fix TLS handshake errors on the client side
		if i%50 == 0 && i != 0 {
			time.Sleep(2 * time.Second)
		}
		wg.Add(1)
		//get responses concurrently with go routines
		go getSpeech(eCtx, newCtx, cancel, nextText, &wg, basepath, tag)
	}
	wg.Wait()

	if newCtx.Err() != nil {
		eCtx.Logger().Error(newCtx.Err())
		return newCtx.Err()
	}

	return nil
}

func getSpeech(
	eCtx echo.Context,
	ctx context.Context,
	cancel context.CancelFunc,
	translate db.Translate,
	wg *sync.WaitGroup,
	basepath,
	tag string) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return // Error somewhere, terminate
	default:
		client, err := texttospeech.NewClient(ctx)
		if err != nil {
			eCtx.Logger().Error(fmt.Errorf("error creating texttospeech client: %s", err))
			cancel()
			return
		}
		defer client.Close()

		// Perform the text-to-speech request on the text input with the selected
		// voice parameters and audio file type.
		req := texttospeechpb.SynthesizeSpeechRequest{
			// Set the text input to be synthesized.
			Input: &texttospeechpb.SynthesisInput{
				InputSource: &texttospeechpb.SynthesisInput_Text{Text: translate.Phrase},
			},
			// Build the voice request, select the language code ("en-US") and the SSML
			// voice gender ("neutral").
			Voice: &texttospeechpb.VoiceSelectionParams{
				LanguageCode: tag,
				SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
				//Name: "af-ZA-Standard-A",
			},
			// Select the type of audio file you want returned.
			AudioConfig: &texttospeechpb.AudioConfig{
				AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			},
		}

		resp, err := client.SynthesizeSpeech(ctx, &req)
		if err != nil {
			eCtx.Logger().Error(fmt.Errorf("error creating Synthesize Speech client: %s", err))
			cancel()
			return
		}

		// The resp AudioContent is binary.
		filename := basepath + strconv.FormatInt(translate.PhraseID, 10)
		err = os.WriteFile(filename, resp.AudioContent, 0644)
		if err != nil {
			eCtx.Logger().Error(fmt.Errorf("error creating translate client: %s", err))
			cancel()
			return
		}
	}
}

func (t *Translate) TranslatePhrases(eCtx echo.Context, ts []db.Translate, dbLang db.Language) ([]util.TranslatesReturn, error) {

	// get language tag to translate to
	langTag, err := language.Parse(dbLang.Tag)
	if err != nil {
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
		go getTranslate(eCtx, newCtx, cancel, langTag, nextTranslate, responses, i, &wg)
	}
	wg.Wait()

	if newCtx.Err() != nil {
		eCtx.Logger().Error(newCtx.Err())
		return nil, newCtx.Err()
	}

	return responses, nil
}

func getTranslate(eCtx echo.Context,
	ctx context.Context,
	cancel context.CancelFunc,
	lang language.Tag,
	phrase db.Translate,
	responses []util.TranslatesReturn,
	i int,
	wg *sync.WaitGroup) {

	defer wg.Done()
	select {
	case <-ctx.Done():
		return // Error somewhere, terminate
	default: // Default to avoid blocking
		client, err := translate.NewClient(ctx)
		if err != nil {
			eCtx.Logger().Error(fmt.Errorf("error creating translate client: %s", err))
			cancel()
			return
		}
		defer client.Close()

		resp, err := client.Translate(ctx, []string{phrase.Phrase}, lang, nil)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				return
			default:
				eCtx.Logger().Error(fmt.Errorf("error translating text: %s", err))
				cancel()
			}
			return
		}

		if len(resp) == 0 {
			eCtx.Logger().Error(fmt.Errorf("translate returned empty response to text: %s", err))
			cancel()
		}

		//app.Logger.PrintInfo(fmt.Sprintf("response is: %s", resp[0].Text), nil)
		responses[i] = util.TranslatesReturn{
			PhraseId: phrase.PhraseID,
			Text:     resp[0].Text,
		}
	}
}

func (t *Translate) CreateTTS(eCtx echo.Context, q db.Querier, lang db.Language, title db.Title, basepath string) error {
	skip, err := pathExists(basepath)
	if err != nil {
		eCtx.Logger().Error(err)
		return err
	}

	if !skip {
		fromTranslates, err := getOrCreateTranslates(eCtx, q, t, title.ID, lang, title.OgLanguageID)
		if err != nil {
			return err
		}

		err = os.MkdirAll(basepath, 0777)
		if err != nil {
			eCtx.Logger().Error(err)
			return err
		}

		if err = t.TextToSpeech(eCtx, fromTranslates, basepath, lang.Tag); err != nil {
			eCtx.Logger().Error(err)
			return err
		}
	}

	return nil
}

func getOrCreateTranslates(eCtx echo.Context, q db.Querier, t TranslateX, titleId int64, toLang db.Language, fromLangId int16) ([]db.Translate, error) {
	// see if translates exist for title for language
	exists, err := q.SelectExistsTranslates(
		eCtx.Request().Context(),
		db.SelectExistsTranslatesParams{
			LanguageID: toLang.ID,
			ID:         titleId,
		})

	// if exists get translates for language
	if exists {
		params := db.SelectTranslatesByTitleIdLangIdParams{
			LanguageID: toLang.ID,
			ID:         titleId,
		}
		translates, err := q.SelectTranslatesByTitleIdLangId(eCtx.Request().Context(), params)
		if err != nil {
			eCtx.Logger().Error(err)
			return nil, err
		}
		return translates, nil
	}

	// if not exists get translates for fromLangId
	fromTranslates, err := q.SelectTranslatesByTitleIdLangId(
		eCtx.Request().Context(),
		db.SelectTranslatesByTitleIdLangIdParams{
			LanguageID: fromLangId,
			ID:         titleId,
		})
	if err != nil {
		eCtx.Logger().Error(err)
		return nil, err
	}

	// create translates for title and to language and return
	translatesReturn, err := t.TranslatePhrases(eCtx, fromTranslates, toLang)
	if err != nil {
		eCtx.Logger().Error(err)
		return nil, err
	}

	dbTranslates, err := t.InsertTranslates(eCtx, q, toLang.ID, translatesReturn)
	if err != nil {
		eCtx.Logger().Error(err)
		return nil, err
	}
	return dbTranslates, nil
}

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

// pathExists returns whether the given file or directory exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
