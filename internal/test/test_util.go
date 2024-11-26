package test

import (
	"fmt"
	"math"
	"math/rand"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/oapi"
)

var (
	AudioBasePath = GetProjectRoot() + "./../tmp/test/audio/"
)

func GetProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(filename))
}

func RequireMatchAnyExcept(t *testing.T, model any, response any, skip []string, except string, shouldEqual any) {
	v := reflect.ValueOf(response)
	u := reflect.ValueOf(model)

	for i := 0; i < v.NumField(); i++ {
		// Check if field name is the one that should be different
		if v.Type().Field(i).Name == except {
			// Check if type is int32 or int64
			if v.Field(i).CanInt() {
				// check if equal as int64
				require.Equal(t, shouldEqual, v.Field(i).Int())
			} else {
				// if not check if equal as string
				require.Equal(t, shouldEqual, v.Field(i).String())
			}
		} else if slices.Contains(skip, v.Type().Field(i).Name) {
			continue
		} else {
			if v.Field(i).CanInt() {
				require.Equal(t, u.Field(i).Int(), v.Field(i).Int())
			} else {
				require.Equal(t, u.Field(i).String(), v.Field(i).String())
			}
		}
	}
}

const (
	ValidTitleId       = -1
	ValidOgLanguageId  = -1
	ValidNewLanguageId = -1
	alphabet           = "abcdefghijklmnopqrstuvwxyz"
)

// RandomInt64 generates a random integer between min and max
func RandomInt64() int64 {
	return rand.Int63n(math.MaxInt64 - 1) //nolint:gosec
}

// RandomInt32 generates a random integer between min and max
func RandomInt32() int32 {
	return rand.Int31n(math.MaxInt32 - 1) //nolint:gosec
}

// RandomInt16 generates a random integer between min and max
func RandomInt16() int16 {
	return int16(rand.Int()) //nolint:gosec
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)] //nolint:gosec
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomPhrase() oapi.Phrase {
	return oapi.Phrase{
		Id:      RandomInt64(),
		TitleId: RandomInt64(),
	}
}

// RandomVoice creates a random db Voice for testing
func RandomVoice() (voice db.Voice) {
	return db.Voice{
		ID:                     RandomInt16(),
		LanguageID:             RandomInt16(),
		LanguageCodes:          []string{RandomString(8), RandomString(8)},
		SsmlGender:             "FEMALE",
		Name:                   RandomString(8),
		NaturalSampleRateHertz: 24000,
	}
}

func RandomTitle() (title db.Title) {
	return db.Title{
		ID:           RandomInt64(),
		Title:        RandomString(8),
		NumSubs:      RandomInt16(),
		OgLanguageID: ValidOgLanguageId,
	}
}

// RandomTranslate create a random db Translate for testing
func RandomTranslate(phrase oapi.Phrase, languageId int16) db.Translate {
	return db.Translate{
		PhraseID:   phrase.Id,
		LanguageID: languageId,
		Phrase:     RandomString(8) + " " + RandomString(8),
		PhraseHint: RandomString(8) + " " + RandomString(8),
	}
}
