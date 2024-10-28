package api

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"strconv"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/oapi"
	"talkliketv.click/tltv/internal/test"
	"testing"
)

func TestGetPhrases(t *testing.T) {
	user, _ := randomUser(t)
	phrase := test.RandomPhrase()
	ogTranslate := randomTranslate(phrase, user.OgLanguageID)
	newTranslate := randomTranslate(phrase, user.NewLanguageID)

	selectPhrasesFromTranslatesParams := db.SelectPhrasesFromTranslatesWithCorrectParams{
		LanguageID:   user.OgLanguageID,
		LanguageID_2: user.OgLanguageID,
		UserID:       user.ID,
		TitleID:      user.TitleID,
		Limit:        10,
	}

	selectPhrasesFromTranslatesRow := db.SelectPhrasesFromTranslatesWithCorrectRow{
		PhraseID:     phrase.Id,
		Phrase:       ogTranslate.Phrase,
		PhraseHint:   ogTranslate.PhraseHint,
		Phrase_2:     newTranslate.Phrase,
		PhraseHint_2: newTranslate.PhraseHint,
	}

	selectPhrasesFromTranslatesRowList := []db.SelectPhrasesFromTranslatesWithCorrectRow{selectPhrasesFromTranslatesRow}

	testCases := []testCase{
		{
			name:   "OK",
			user:   user,
			values: map[string]any{"limit": true},
			buildStubs: func(stubs buildStubs) {
				stubs.mockQuerier.EXPECT().
					SelectPhrasesFromTranslatesWithCorrect(gomock.Any(), selectPhrasesFromTranslatesParams).
					Times(1).
					Return(selectPhrasesFromTranslatesRowList, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var got []db.SelectPhrasesFromTranslatesRow
				err := json.Unmarshal([]byte(body), &got)
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, selectPhrasesFromTranslatesRowList[0], got[0], nil, "", "")
			},
			permissions: []string{db.ReadTitlesCode},
		},
		{
			name:   "No limit set",
			user:   user,
			values: map[string]any{"limit": false},
			buildStubs: func(stubs buildStubs) {
				stubs.mockQuerier.EXPECT().
					SelectPhrasesFromTranslatesWithCorrect(gomock.Any(), selectPhrasesFromTranslatesParams).
					Times(1).
					Return(selectPhrasesFromTranslatesRowList, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var got []db.SelectPhrasesFromTranslatesWithCorrectRow
				err := json.Unmarshal([]byte(body), &got)
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, selectPhrasesFromTranslatesRowList[0], got[0], nil, "", "")
			},
			permissions: []string{db.ReadTitlesCode},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts, jwsToken := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, []byte(""), ts, phrasesBasePath, http.MethodGet, jwsToken)

			// add parameters to the url query path
			q := req.URL.Query()
			if tc.values["limit"] == true {
				q.Add("limit", "10")
			}
			req.URL.RawQuery = q.Encode()

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}

func TestUpdateUsersPhrases(t *testing.T) {
	user1, _ := randomUser(t)
	phrase := test.RandomPhrase()
	usersPhrase := randomUsersPhrase(user1, phrase)

	updateUsersPhrasesParams := db.UpdateUsersPhrasesByThreeIdsParams{
		TitleID:       usersPhrase.TitleID,
		LanguageID:    usersPhrase.LanguageID,
		UserID:        usersPhrase.UserID,
		PhraseCorrect: usersPhrase.PhraseCorrect,
		PhraseID:      phrase.Id,
	}

	testCases := []testCase{
		{
			name:        "update UsersPhrase phraseCorrect",
			permissions: []string{db.ReadTitlesCode},
			values: map[string]any{
				"phraseId":   strconv.FormatInt(phrase.Id, 10),
				"languageId": fmt.Sprint(user1.NewLanguageID)},
			user:     user1,
			extraInt: user1.ID,
			body: `[
			{
				"op": "replace",
				"path": "/phraseCorrect",
				"value": 1
			}
		]`,
			buildStubs: func(stubs buildStubs) {
				args := db.SelectUsersPhrasesByIdsParams{
					UserID:     user1.ID,
					LanguageID: user1.NewLanguageID,
					PhraseID:   phrase.Id,
				}
				usersPhraseCopy := usersPhrase
				paramsCopy := updateUsersPhrasesParams
				paramsCopy.PhraseCorrect = 1
				usersPhraseCopy.PhraseCorrect = 1
				stubs.mockQuerier.EXPECT().
					SelectUsersPhrasesByIds(gomock.Any(), args).
					Times(1).
					Return(usersPhrase, nil)
				stubs.mockQuerier.EXPECT().
					UpdateUsersPhrasesByThreeIds(gomock.Any(), paramsCopy).
					Times(1).
					Return(usersPhraseCopy, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var got db.UsersPhrase
				err := json.Unmarshal([]byte(body), &got)
				require.NoError(t, err)
				var x int64 = 1
				test.RequireMatchAnyExcept(t, usersPhrase, got, []string{}, "PhraseCorrect", x)
			},
		},
		{
			name:        "invalid format of patch",
			permissions: []string{db.ReadTitlesCode},
			values: map[string]any{
				"phraseId":   strconv.FormatInt(phrase.Id, 10),
				"languageId": fmt.Sprint(user1.NewLanguageID)},
			user:     user1,
			extraInt: user1.ID,
			body: `[
			{
				"wrong": "replace",
				"path": "/phraseCorrect",
				"value": 1
			}
		]`,
			buildStubs: func(stubs buildStubs) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "\"message\":\"Error at \\\"/0\\\": property \\\"wrong\\\" is unsupported\\nSchema:\\n ")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlPath := usersPhrasesBasePath + "/" + tc.values["phraseId"].(string) + "/" + tc.values["languageId"].(string)
			ts, jwsToken := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, []byte(tc.body.(string)), ts, urlPath, http.MethodPatch, jwsToken)

			// change request content-type to patch
			req.Header.Set("Content-Type", "application/json-patch+json")
			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}

func randomUsersPhrase(user db.User, phrase oapi.Phrase) db.UsersPhrase {
	return db.UsersPhrase{
		PhraseID:      phrase.Id,
		TitleID:       phrase.TitleId,
		UserID:        user.ID,
		LanguageID:    user.OgLanguageID,
		PhraseCorrect: 0,
	}
}
