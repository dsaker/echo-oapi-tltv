package api

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	mockdb "talkliketv.click/tltv/db/mock"
	db "talkliketv.click/tltv/db/sqlc"
	mock "talkliketv.click/tltv/internal/mock"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

type translatesTestCase struct {
	name              string
	buildStubs        func(*mockdb.MockQuerier, *mock.MockTranslateX)
	checkTranslate    func([]db.Translate, error)
	checkTranslateRow func([]util.TranslatesReturn, error)
}

func TestInsertNewPhrases(t *testing.T) {
	title := randomTitle()
	title.OgLanguageID = 27
	randomPhrase1 := randomPhrase()
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
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
				store.EXPECT().
					InsertPhrases(gomock.Any(), title.ID).
					Times(1).
					Return(dbPhrase1, nil)
				store.EXPECT().
					InsertTranslates(gomock.Any(), insertTranslatesParams).
					Times(1).
					Return(translate1, nil)
			},
			checkTranslate: func(translates []db.Translate, err error) {
				require.NoError(t, err)
				require.Contains(t, translates, translate1)
				requireMatchAnyExcept(t, translates[0], translate1, nil, "", "")
			},
		},
		{
			name: "DB Connection Error",
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
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

			text := mock.NewMockTranslateX(ctrl)
			store := mockdb.NewMockQuerier(ctrl)
			tc.buildStubs(store, text)

			e, _ := NewServer(testCfg, store, text)

			req := httptest.NewRequest(http.MethodPost, "/titles/translates", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			translate := Translate{}
			translates, err := translate.InsertNewPhrases(c, title, store, stringsSlice)
			tc.checkTranslate(translates, err)
		})
	}
}

func TestInsertTranslates(t *testing.T) {
	title := randomTitle()
	title.OgLanguageID = 27
	newLanguageId := 109
	randomPhrase1 := randomPhrase()
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
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
				store.EXPECT().
					InsertTranslates(gomock.Any(), insertTranslatesParams).
					Times(1).
					Return(translate1, nil)
			},
			checkTranslate: func(translates []db.Translate, err error) {
				require.NoError(t, err)
				require.Contains(t, translates, translate1)
				requireMatchAnyExcept(t, translates[0], translate1, nil, "", "")
			},
		},
		{
			name: "DB Connection Error",
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
				store.EXPECT().
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

			text := mock.NewMockTranslateX(ctrl)
			store := mockdb.NewMockQuerier(ctrl)
			tc.buildStubs(store, text)

			e, _ := NewServer(testCfg, store, text)

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
	title := randomTitle()
	title.OgLanguageID = 27

	basepath := "/tmp/" + strconv.FormatInt(title.ID, 10) + "/"
	err := os.MkdirAll(basepath, 0777)
	require.NoError(t, err)
	defer os.RemoveAll(basepath)

	newLanguage := language.English
	randomPhrase1 := randomPhrase()
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
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
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

			text := mock.NewMockTranslateX(ctrl)
			store := mockdb.NewMockQuerier(ctrl)
			tc.buildStubs(store, text)

			e, _ := NewServer(testCfg, store, text)
			req := httptest.NewRequest(http.MethodPost, "/titles/translates", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			translate := Translate{}
			err = translate.TextToSpeech(c, []db.Translate{translate1}, basepath, newLanguage.String())
			tc.checkTranslate(nil, err)
		})
	}
}

func TestTranslatePhrases(t *testing.T) {
	title := randomTitle()
	title.OgLanguageID = 27

	newLanguage := language.Spanish
	randomPhrase1 := randomPhrase()
	text1 := "This is sentence one."
	translate1 := db.SelectTranslatesByTitleIdLangIdRow{
		PhraseID: randomPhrase1.Id,
		Phrase:   text1,
	}

	translatesReturn := []util.TranslatesReturn{{PhraseId: 0, Text: "Esta es la primera oración."}}

	testCases := []translatesTestCase{
		{
			name: "No error",
			buildStubs: func(store *mockdb.MockQuerier, text *mock.MockTranslateX) {
			},
			checkTranslateRow: func(translatesRow []util.TranslatesReturn, err error) {
				require.NoError(t, err)
				requireMatchAnyExcept(t, translatesRow[0], translatesReturn[0], nil, "PhraseId", translatesReturn[0].PhraseId)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			text := mock.NewMockTranslateX(ctrl)
			store := mockdb.NewMockQuerier(ctrl)
			tc.buildStubs(store, text)

			e, _ := NewServer(testCfg, store, text)
			req := httptest.NewRequest(http.MethodPost, "/titles/translates", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			translate := Translate{}
			translatesRow, err := translate.TranslatePhrases(c, []db.SelectTranslatesByTitleIdLangIdRow{translate1}, newLanguage)
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
