package translates

import (
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"cloud.google.com/go/translate"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	mockdb "talkliketv.click/tltv/db/mock"
	db "talkliketv.click/tltv/db/sqlc"
	mockc "talkliketv.click/tltv/internal/mock/clients"
	mockt "talkliketv.click/tltv/internal/mock/translates"
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/test"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

type translatesTestCase struct {
	name              string
	buildStubs        func(*mockdb.MockQuerier, *mockt.MockTranslateX, *mockc.MockTranslateClientX, *mockc.MockTTSClientX)
	checkTranslate    func([]db.Translate, error)
	checkTranslateRow func([]util.TranslatesReturn, error)
}

func TestInsertNewPhrases(t *testing.T) {
	title := RandomTitle()
	title.OgLanguageID = 27
	randomPhrase1 := test.RandomPhrase()
	text1 := "This is sentence one."
	hintString1 := makeHintString(text1)
	translate1 := db.Translate{
		PhraseID:   randomPhrase1.Id,
		LanguageID: title.OgLanguageID,
		Phrase:     text1,
		PhraseHint: hintString1,
	}

	dbPhrase1 := db.Phrase{
		ID:      randomPhrase1.Id,
		TitleID: title.ID,
	}

	stringsSlice := []string{text1}

	insertTranslatesParams := db.InsertTranslatesParams{
		PhraseID:   randomPhrase1.Id,
		LanguageID: title.OgLanguageID,
		Phrase:     text1,
		PhraseHint: hintString1,
	}

	testCases := []translatesTestCase{
		{
			name: "No error",
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX, tc *mockc.MockTranslateClientX, tts *mockc.MockTTSClientX) {
				//InsertNewPhrases(e echo.Context, title db.Title, q db.Querier, stringsSlice []string) ([]db.Translate, error)
				store.EXPECT().InsertPhrases(gomock.Any(), title.ID).
					Return(dbPhrase1, nil)
				store.EXPECT().InsertTranslates(gomock.Any(), insertTranslatesParams).
					Return(translate1, nil)
			},
			checkTranslate: func(translates []db.Translate, err error) {
				require.NoError(t, err)
				require.Contains(t, translates, translate1)
				test.RequireMatchAnyExcept(t, translates[0], translate1, nil, "", "")
			},
		},
		{
			name: "DB Connection Error",
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX, tc *mockc.MockTranslateClientX, tts *mockc.MockTTSClientX) {
				store.EXPECT().
					InsertPhrases(gomock.Any(), title.ID).
					Times(1).
					Return(db.Phrase{}, sql.ErrConnDone)
			},
			checkTranslate: func(translates []db.Translate, err error) {
				require.Contains(t, err.Error(), "sql: connection is already closed")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			text := mockt.NewMockTranslateX(ctrl)
			store := mockdb.NewMockQuerier(ctrl)
			tclient := mockc.NewMockTranslateClientX(ctrl)
			ttsclient := mockc.NewMockTTSClientX(ctrl)
			tc.buildStubs(store, text, tclient, ttsclient)

			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/titles/translates", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			translate := &Translate{}
			translates, err := translate.InsertNewPhrases(c, title, store, stringsSlice)
			tc.checkTranslate(translates, err)
		})
	}
}

func TestInsertTranslates(t *testing.T) {
	title := RandomTitle()
	title.OgLanguageID = 27
	newLanguageId := 109
	randomPhrase1 := RandomPhrase()
	text1 := "This is sentence one."
	hintString1 := makeHintString(text1)
	translate1 := db.Translate{
		PhraseID:   randomPhrase1.Id,
		LanguageID: title.OgLanguageID,
		Phrase:     text1,
		PhraseHint: hintString1,
	}

	translatesReturn := util.TranslatesReturn{
		PhraseId: randomPhrase1.Id,
		Text:     text1,
	}

	insertTranslatesParams := db.InsertTranslatesParams{
		PhraseID:   randomPhrase1.Id,
		LanguageID: int16(newLanguageId),
		Phrase:     text1,
		PhraseHint: hintString1,
	}

	testCases := []translatesTestCase{
		{
			name: "No error",
			buildStubs: func(s *mockdb.MockQuerier, t *mockt.MockTranslateX, tc *mockc.MockTranslateClientX, tts *mockc.MockTTSClientX) {
				s.EXPECT().
					InsertTranslates(gomock.Any(), insertTranslatesParams).
					Times(1).
					Return(translate1, nil)
			},
			checkTranslate: func(translates []db.Translate, err error) {
				require.NoError(t, err)
				require.Contains(t, translates, translate1)
				test.RequireMatchAnyExcept(t, translates[0], translate1, nil, "", "")
			},
		},
		{
			name: "DB Connection Error",
			buildStubs: func(s *mockdb.MockQuerier, t *mockt.MockTranslateX, tc *mockc.MockTranslateClientX, tts *mockc.MockTTSClientX) {
				s.EXPECT().
					InsertTranslates(gomock.Any(), insertTranslatesParams).
					Times(1).
					Return(db.Translate{}, sql.ErrConnDone)
			},
			checkTranslate: func(translates []db.Translate, err error) {
				require.Contains(t, err.Error(), "sql: connection is already closed")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			text := mockt.NewMockTranslateX(ctrl)
			store := mockdb.NewMockQuerier(ctrl)
			tclient := mockc.NewMockTranslateClientX(ctrl)
			ttsclient := mockc.NewMockTTSClientX(ctrl)
			tc.buildStubs(store, text, tclient, ttsclient)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/titles/translates", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			translate := Translate{}
			translates, err := translate.InsertTranslates(c, store, int16(newLanguageId), []util.TranslatesReturn{translatesReturn})
			tc.checkTranslate(translates, err)
		})
	}
}

func TestTextToSpeech(t *testing.T) {
	title := RandomTitle()
	title.OgLanguageID = 27

	basepath := "/tmp/" + strconv.FormatInt(title.ID, 10) + "/"
	err := os.MkdirAll(basepath, 0777)
	require.NoError(t, err)
	defer os.RemoveAll(basepath)

	newLanguage := language.English
	randomPhrase1 := RandomPhrase()
	text1 := "This is sentence one."
	hintString1 := makeHintString(text1)
	translate1 := db.Translate{
		PhraseID:   randomPhrase1.Id,
		LanguageID: title.OgLanguageID,
		Phrase:     text1,
		PhraseHint: hintString1,
	}

	testCases := []translatesTestCase{
		{
			name: "No error",
			buildStubs: func(s *mockdb.MockQuerier, t *mockt.MockTranslateX, tc *mockc.MockTranslateClientX, tts *mockc.MockTTSClientX) {
				req := texttospeechpb.SynthesizeSpeechRequest{
					// Set the text input to be synthesized.
					Input: &texttospeechpb.SynthesisInput{
						InputSource: &texttospeechpb.SynthesisInput_Text{Text: text1},
					},
					// Build the voice request, select the language code ("en-US") and the SSML
					// voice gender ("neutral").
					Voice: &texttospeechpb.VoiceSelectionParams{
						LanguageCode: language.English.String(),
						SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
						//Name: "af-ZA-Standard-A",
					},
					// Select the type of audio file you want returned.
					AudioConfig: &texttospeechpb.AudioConfig{
						AudioEncoding: texttospeechpb.AudioEncoding_MP3,
					},
				}
				resp := texttospeechpb.SynthesizeSpeechResponse{}
				tts.EXPECT().SynthesizeSpeech(gomock.Any(), &req).Return(&resp, nil)
			},
			checkTranslate: func(translates []db.Translate, err error) {
				require.NoError(t, err)
				isEmpty, err := IsDirectoryEmpty(basepath)
				require.NoError(t, err)
				require.False(t, isEmpty)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			text := mockt.NewMockTranslateX(ctrl)
			store := mockdb.NewMockQuerier(ctrl)
			trc := mockc.NewMockTranslateClientX(ctrl)
			tts := mockc.NewMockTTSClientX(ctrl)
			tc.buildStubs(store, text, trc, tts)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/titles/translates", nil)
			rec := httptest.NewRecorder()
			newE := e.NewContext(req, rec)

			translates := &Translate{}
			err = translates.TextToSpeech(newE, []db.Translate{translate1}, tts, basepath, newLanguage.String())
			tc.checkTranslate(nil, err)
		})
	}
}

func TestTranslatePhrases(t *testing.T) {
	title := RandomTitle()
	title.OgLanguageID = 27

	newLanguage := db.Language{
		ID:       109,
		Language: "Spanish",
		Tag:      "es",
	}
	randomPhrase1 := RandomPhrase()
	text1 := "This is sentence one."
	translate1 := db.Translate{
		PhraseID: randomPhrase1.Id,
		Phrase:   text1,
	}

	translatesReturn := []util.TranslatesReturn{{PhraseId: 0, Text: "Esta es la primera oración."}}

	translation := translate.Translation{Text: "Esta es la primera oración."}
	testCases := []translatesTestCase{
		{
			name: "No error",
			buildStubs: func(s *mockdb.MockQuerier, t *mockt.MockTranslateX, tr *mockc.MockTranslateClientX, ts *mockc.MockTTSClientX) {
				tr.EXPECT().Translate(gomock.Any(), []string{text1}, language.Spanish, nil).
					Return([]translate.Translation{translation}, nil)
			},
			checkTranslateRow: func(translatesRow []util.TranslatesReturn, err error) {
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, translatesRow[0], translatesReturn[0], nil, "PhraseId", translatesReturn[0].PhraseId)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			text := mockt.NewMockTranslateX(ctrl)
			store := mockdb.NewMockQuerier(ctrl)
			client := mockc.NewMockTranslateClientX(ctrl)
			tts := mockc.NewMockTTSClientX(ctrl)
			tc.buildStubs(store, text, client, tts)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/titles/translates", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			translate := Translate{}
			translatesRow, err := translate.TranslatePhrases(c, []db.Translate{translate1}, newLanguage, client)
			tc.checkTranslateRow(translatesRow, err)
		})
	}
}

// IsDirectoryEmpty returns true if directory is empty and false if not
func IsDirectoryEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Read only one entry
	if err == io.EOF {
		return true, nil // Directory is empty
	}
	return false, nil // Directory is not empty
}

func RandomPhrase() oapi.Phrase {
	return oapi.Phrase{
		Id:      test.RandomInt64(),
		TitleId: test.RandomInt64(),
	}
}

func RandomTitle() (title db.Title) {

	return db.Title{
		ID:           test.RandomInt64(),
		Title:        test.RandomString(8),
		NumSubs:      test.RandomInt16(),
		OgLanguageID: test.ValidOgLanguageId,
	}
}