package api

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	mockdb "talkliketv.click/tltv/db/mock"
	"testing"
)

func TestGetTitle(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		body          map[string]any
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(res *http.Response)
	}{
		{
			name: "Invalid Username",
			body: map[string]any{
				"name":          "invalid-user#1",
				"email":         user.Email,
				"flipped":       user.Flipped,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      password,
				"titleId":       user.TitleID,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "InvalidEmail",
			body: map[string]any{
				"flipped":       user.Flipped,
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      password,
				"titleId":       user.TitleID,
				"email":         "invalid-email",
			},
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "string doesn't match the regular expression ")
			},
		},
		{
			name: "TooShortPassword",
			body: map[string]any{
				"flipped":       user.Flipped,
				"name":          user.Name,
				"newLanguageId": user.NewLanguageID,
				"ogLanguageId":  user.OgLanguageID,
				"password":      "string",
				"titleId":       user.TitleID,
				"email":         user.Email,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				body := readBody(t, res)
				require.Contains(t, body, "minimum string length is 8")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockQuerier(ctrl)
			tc.buildStubs(store)

			spec, err := GetSwagger()
			require.NoError(t, err)

			e := echo.New()
			_ = NewServer(e, testCfg, store, spec)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			ts := httptest.NewServer(e)
			urlPath := "/v1/users"
			res, err := ts.Client().Post(ts.URL+urlPath, "application/json", bytes.NewReader(data))
			require.NoError(t, err)
			tc.checkResponse(res)
		})
	}
}
