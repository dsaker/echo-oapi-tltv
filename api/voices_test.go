package api

import (
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/test"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

func TestListVoices(t *testing.T) {
	user, _ := randomUser(t)

	voice1 := util.RandomVoice()
	voice2 := util.RandomVoice()
	voices := []db.Voice{voice1, voice2}

	testCases := []testCase{
		{
			name: "OK",
			user: user,
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					ListVoices(gomock.Any()).
					Times(1).
					Return(voices, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var gotVoices []db.Voice
				err := json.Unmarshal([]byte(body), &gotVoices)
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, gotVoices[0], voices[0], nil, "", "")
			},
			permissions: []string{},
		},
		{
			name:   "With language id",
			user:   user,
			values: map[string]any{"languageId": true},
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					SelectVoicesByLanguageId(gomock.Any(), int16(1)).
					Times(1).
					Return(voices, nil)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				body := readBody(t, res)
				var gotVoices []db.Voice
				err := json.Unmarshal([]byte(body), &gotVoices)
				require.NoError(t, err)
				test.RequireMatchAnyExcept(t, gotVoices[0], voices[0], nil, "", "")
			},
			permissions: []string{},
		},
		{
			name: "conn already closed",
			user: user,
			buildStubs: func(stubs MockStubs) {
				stubs.MockQuerier.EXPECT().
					ListVoices(gomock.Any()).
					Times(1).
					Return([]db.Voice{}, sql.ErrConnDone)
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusInternalServerError, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "sql: connection is already closed")
			},
			permissions: []string{},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts, jwsToken := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, []byte(""), ts, voicesBasePath, http.MethodGet, jwsToken)
			q := req.URL.Query()
			if tc.values["languageId"] == true {
				q.Add("languageId", "1")
			}
			req.URL.RawQuery = q.Encode()

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}
