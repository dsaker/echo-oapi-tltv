package api

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/test"
	"testing"
)

func TestListLanguages(t *testing.T) {
	user, _ := randomUser(t)

	language1 := randomLanguage()
	language2 := randomLanguage()
	languages := []db.Language{language1, language2}

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			buildStubs: func(stubs buildStubs) {
				stubs.mdb.EXPECT().
					ListLanguages(gomock.Any()).
					Times(1).
					Return(languages, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var gotLanguages []db.Language
				err := json.Unmarshal([]byte(body), &gotLanguages)
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, gotLanguages[0], languages[0], nil, "", "")
			},
			permissions: []string{db.ReadTitlesCode},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			//req, ts := setupServerTest(t, ctrl, tc, []byte(""), titlesBasePath, http.MethodGet)
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

			tc.checkResponse(res)
		})
	}
}
