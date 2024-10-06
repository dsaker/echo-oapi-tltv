package api

import (
	"bufio"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"cloud.google.com/go/translate"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	db "talkliketv.click/tltv/db/sqlc"
	"time"
	"unicode"
)

// FindTitles implements all the handlers in the ServerInterface
func (s *Server) FindTitles(ctx echo.Context, params FindTitlesParams) error {

	titles, err := s.queries.ListTitles(
		ctx.Request().Context(),
		db.ListTitlesParams{
			Similarity: params.Similarity,
			Limit:      params.Limit,
		})

	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, titles)
}

func (s *Server) AddTitle(eCtx echo.Context) error {

	lang := eCtx.FormValue("languageId")
	titleName := eCtx.FormValue("titleName")
	langIdInt16, err := strconv.ParseInt(lang, 10, 16)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	//Get language model from id for tag
	langModel, err := s.queries.SelectLanguagesById(eCtx.Request().Context(), int16(langIdInt16))
	if err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	// Get handler for filename, size and headers
	file, err := eCtx.FormFile("filePath")
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	// Check if file size is too large
	if file.Size > 32768 {
		return eCtx.String(http.StatusBadRequest, "file too large")
	}
	src, err := file.Open()
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	eCtx.Logger().Info(fmt.Sprintf("File uploaded successfully: %s", file.Filename))

	// Create strings slice and count number of lines form titles model
	scanner := bufio.NewScanner(src)
	var stringsSlice []string
	numLines := 0
	for scanner.Scan() {
		numLines += 1
		stringsSlice = append(stringsSlice, scanner.Text())
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	title, err := s.queries.InsertTitle(
		eCtx.Request().Context(),
		db.InsertTitleParams{
			Title:        titleName,
			NumSubs:      int16(numLines),
			OgLanguageID: int16(langIdInt16),
		})

	// insert phrases into db as translates object of OgLanguage
	translatesSlice, err := insertPhrases(eCtx, title, s.queries, stringsSlice, numLines)
	if err != nil {
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	//create base path for storing mp3 audio files
	audioBasePath := s.config.TTSBasePath +
		strconv.FormatInt(title.ID, 10) + "/" +
		strconv.Itoa(int(title.OgLanguageID)) + "/"
	err = os.MkdirAll(audioBasePath, os.ModePerm)
	if err != nil {
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}
	err = textToSpeech(eCtx, translatesSlice, audioBasePath, langModel.Tag)

	if err != nil {
		eCtx.Logger().Error(err)
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusOK, title)
}

func (s *Server) FindTitleByID(ctx echo.Context, id int64) error {

	title, err := s.queries.SelectTitleById(ctx.Request().Context(), id)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, title)
}

func (s *Server) DeleteTitle(ctx echo.Context, id int64) error {

	err := s.queries.DeleteTitleById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (s *Server) TranslateTitle(eCtx echo.Context) error {

	var newTranslateTitle TranslateTitleJSONBody
	err := eCtx.Bind(&newTranslateTitle)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	exists, err := s.queries.SelectExistsTranslates(
		eCtx.Request().Context(),
		db.SelectExistsTranslatesParams{
			LanguageID: newTranslateTitle.NewLanguageId,
			ID:         newTranslateTitle.TitleId,
		})
	if exists {
		return eCtx.String(http.StatusBadRequest, "title already exists in that language")
	}

	title, err := s.queries.SelectTitleById(eCtx.Request().Context(), newTranslateTitle.TitleId)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, "invalid title id")
	}

	dbLang, err := s.queries.SelectLanguagesById(eCtx.Request().Context(), newTranslateTitle.NewLanguageId)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	langTag, err := language.Parse(dbLang.Tag)
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	// We're always asynchronous, so lock unsafe operations below
	s.Lock()
	defer s.Unlock()

	translates, err := s.queries.SelectTranslatesByTitleIdLangId(
		eCtx.Request().Context(),
		db.SelectTranslatesByTitleIdLangIdParams{
			ID:         newTranslateTitle.TitleId,
			LanguageID: title.OgLanguageID,
		})
	if err != nil {
		return eCtx.String(http.StatusBadRequest, err.Error())
	}

	newTranslates, err := translatePhrases(eCtx, int(translates[0].NumSubs), translates, langTag)
	if err != nil {
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	insertTranslates, err := insertPhrases(eCtx, title, s.queries, newTranslates, int(translates[0].NumSubs))
	if err != nil {
		return eCtx.String(http.StatusInternalServerError, err.Error())
	}

	return eCtx.JSON(http.StatusCreated, insertTranslates[0])
}

func insertPhrases(eCtx echo.Context, title db.Title, q db.Querier, stringsSlice []string, numLines int) ([]db.Translate, error) {
	dbTranslates := make([]db.Translate, numLines)
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

func textToSpeech(eCtx echo.Context, translatesSlice []db.Translate, basepath, tag string) error {

	// concurrently get all the audio content from Google text-to-speech
	var wg sync.WaitGroup
	// create context with cancel, so you can cancel all other requests after any error
	newCtx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure it's called to release resources even if no errors

	for i, nextSpeech := range translatesSlice {
		// added intermittent sleep to fix TLS handshake errors on the client side
		if i%50 == 0 && i != 0 {
			time.Sleep(1 * time.Second)
		}
		wg.Add(1)
		//get responses concurrently with go routines
		go getSpeech(eCtx, newCtx, cancel, nextSpeech, &wg, basepath, tag)
	}
	wg.Wait()

	if newCtx.Err() != nil {
		eCtx.Logger().Error(newCtx.Err())
		return eCtx.String(http.StatusInternalServerError, newCtx.Err().Error())
	}

	return nil
}

func getSpeech(eCtx echo.Context,
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

		// The resp's AudioContent is binary.
		filename := basepath + strconv.FormatInt(translate.PhraseID, 10) + ".mp3"
		err = os.WriteFile(filename, resp.AudioContent, 0644)
		if err != nil {
			eCtx.Logger().Error(fmt.Errorf("error creating translate client: %s", err))
			cancel()
			return
		}
		fmt.Printf("Audio content written to file: %v\n", filename)
	}
}

func translatePhrases(eCtx echo.Context, numLines int, translatesSlice []db.SelectTranslatesByTitleIdLangIdRow, tag language.Tag) ([]string, error) {

	// concurrently get all the responses from Google Translate
	var wg sync.WaitGroup
	responses := make([]string, numLines) // string array to hold all the responses
	// create context with cancel, so you can cancel all other requests after any error
	newCtx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure it's called to release resources even if no errors

	for i, nextTranslate := range translatesSlice {
		// added intermittent sleep to fix TLS handshake errors on the client side
		if i%50 == 0 && i != 0 {
			time.Sleep(1 * time.Second)
		}
		wg.Add(1)
		//get responses concurrently with go routines
		go getTranslate(eCtx, newCtx, cancel, tag, nextTranslate.Phrase, responses, i, &wg)
	}
	wg.Wait()

	if newCtx.Err() != nil {
		eCtx.Logger().Error(newCtx.Err())
		return nil, eCtx.String(http.StatusInternalServerError, newCtx.Err().Error())
	}

	return responses, nil
}

func getTranslate(eCtx echo.Context,
	ctx context.Context,
	cancel context.CancelFunc,
	lang language.Tag,
	phrase string,
	responses []string,
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

		resp, err := client.Translate(ctx, []string{phrase}, lang, nil)
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
		responses[i] = resp[0].Text
	}
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

func deleteTitleSendError() {

}
