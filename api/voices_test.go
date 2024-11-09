package api

import (
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/test"
	"testing"
)

func TestListVoices(t *testing.T) {
	user, _ := randomUser(t)

	voice1 := randomVoice()
	voice2 := randomVoice()
	voices := []db.Voice{voice1, voice2}

	testCases := []testCase{
		{
			name:   "OK",
			user:   user,
			values: map[string]any{"language_id": false},
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
			values: map[string]any{"language_id": true},
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
			name:   "conn already closed",
			user:   user,
			values: map[string]any{"language_id": false},
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

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts, jwsToken := setupServerTest(t, ctrl, tc)
			req := jsonRequest(t, []byte(""), ts, voicesBasePath, http.MethodGet, jwsToken)

			// add parameters to the url query path
			q := req.URL.Query()
			if tc.values["language_id"] == true {
				q.Add("language_id", "1")
			}
			req.URL.RawQuery = q.Encode()

			res, err := ts.Client().Do(req)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}
