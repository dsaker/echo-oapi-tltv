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
	db "talkliketv.click/tltv/db/sqlc"
	mockdb "talkliketv.click/tltv/internal/mock/db"
	mockt "talkliketv.click/tltv/internal/mock/translates"
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/test"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

type translatesTestCase struct {
	name              string
	buildStubs        func(*mockdb.MockQuerier, *mockt.MockTranslateX, *mockt.MockTranslateClientX, *mockt.MockTTSClientX)
	checkTranslate    func([]db.Translate, error)
	checkTranslateRow func([]util.TranslatesReturn, error)
}

func TestInsertNewPhrases(t *testing.T) {
	title := RandomTitle()
	title.OgLanguageID = 27
	randomPhrase1 := util.RandomPhrase()
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
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX, tc *mockt.MockTranslateClientX, tts *mockt.MockTTSClientX) {
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
			buildStubs: func(store *mockdb.MockQuerier, text *mockt.MockTranslateX, tc *mockt.MockTranslateClientX, tts *mockt.MockTTSClientX) {
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
			tclient := mockt.NewMockTranslateClientX(ctrl)
			ttsclient := mockt.NewMockTTSClientX(ctrl)
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
			buildStubs: func(s *mockdb.MockQuerier, t *mockt.MockTranslateX, tc *mockt.MockTranslateClientX, tts *mockt.MockTTSClientX) {
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
			buildStubs: func(s *mockdb.MockQuerier, t *mockt.MockTranslateX, tc *mockt.MockTranslateClientX, tts *mockt.MockTTSClientX) {
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
			tclient := mockt.NewMockTranslateClientX(ctrl)
			ttsclient := mockt.NewMockTTSClientX(ctrl)
			tc.buildStubs(store, text, tclient, ttsclient)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/titles/translates", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			translate2 := Translate{}
			dbTranslates, err := translate2.InsertTranslates(c, store, int16(newLanguageId), []util.TranslatesReturn{translatesReturn})
			tc.checkTranslate(dbTranslates, err)
		})
	}
}

func TestTextToSpeech(t *testing.T) {
	title := RandomTitle()
	title.OgLanguageID = 27

	basepath := test.AudioBasePath + strconv.FormatInt(title.ID, 10) + "/"
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
			buildStubs: func(s *mockdb.MockQuerier, t *mockt.MockTranslateX, tc *mockt.MockTranslateClientX, tts *mockt.MockTTSClientX) {
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
			trc := mockt.NewMockTranslateClientX(ctrl)
			tts := mockt.NewMockTTSClientX(ctrl)
			tc.buildStubs(store, text, trc, tts)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/titles/translates", nil)
			rec := httptest.NewRecorder()
			newE := e.NewContext(req, rec)

			translates := New(trc, tts)

			err = translates.TextToSpeech(newE, []db.Translate{translate1}, basepath, newLanguage.String())
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
			buildStubs: func(s *mockdb.MockQuerier, t *mockt.MockTranslateX, tr *mockt.MockTranslateClientX, ts *mockt.MockTTSClientX) {
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
			trc := mockt.NewMockTranslateClientX(ctrl)
			tts := mockt.NewMockTTSClientX(ctrl)
			tc.buildStubs(store, text, trc, tts)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/titles/translates", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			translate2 := New(trc, tts)
			translatesRow, err := translate2.TranslatePhrases(c, []db.Translate{translate1}, newLanguage)
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
		Id:      util.RandomInt64(),
		TitleId: util.RandomInt64(),
	}
}

func RandomTitle() (title db.Title) {

	return db.Title{
		ID:           util.RandomInt64(),
		Title:        util.RandomString(8),
		NumSubs:      util.RandomInt16(),
		OgLanguageID: util.ValidOgLanguageId,
	}
}
