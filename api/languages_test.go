package api

import (
	"encoding/json"
	"net/http"
	"talkliketv.click/tltv/internal/util"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/test"
)

func TestListLanguages(t *testing.T) {
	if util.Integration {
		t.Skip("skipping unit test")
	}

	t.Parallel()

	user, _ := randomUser(t)

	language1 := randomLanguage()
	language2 := randomLanguage()
	languageRow1 := db.ListLanguagesSimilarRow{
		ID:         language1.ID,
		Language:   language1.Language,
		Tag:        language1.Tag,
		Similarity: 0,
	}
	languageRow2 := db.ListLanguagesSimilarRow{
		ID:         language2.ID,
		Language:   language2.Language,
		Tag:        language2.Tag,
		Similarity: 0,
	}
	languages := []db.ListLanguagesSimilarRow{languageRow1, languageRow2}

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					// ListLanguagesSimilar(ctx context.Context, similarity string) ([]ListLanguagesSimilarRow, error)
					ListLanguagesSimilar(gomock.Any(), "").
					Times(1).
					Return(languages, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var gotLanguages []db.ListLanguagesSimilarRow
				err := json.Unmarshal([]byte(body), &gotLanguages)
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, gotLanguages[0], languages[0], nil, "", "")
			},
			permissions: []string{db.ReadTitlesCode},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts, jwsToken := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, []byte(""), ts, languagesBasePath, http.MethodGet, jwsToken)
			q := req.URL.Query()
			if tc.values["similarity"] == true {
				q.Add("similarity", "similar")
			}
			if tc.values["limit"] == true {
				q.Add("limit", "10")
			}
			req.URL.RawQuery = q.Encode()

			res, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			tc.checkResponse(res)
		})
	}
}
