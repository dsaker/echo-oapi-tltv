package util

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/oapi"
)

const (
	ValidTitleId       = -1
	ValidOgLanguageId  = -1
	ValidNewLanguageId = -1
	ValidPermissionId  = 1
	InvalidUserId      = -2
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomInt64 generates a random integer between min and max
func RandomInt64() int64 {
	return rand.Int63n(math.MaxInt64 - 1)
}

// RandomInt32 generates a random integer between min and max
func RandomInt32() int32 {
	return rand.Int31n(math.MaxInt32 - 1)
}

// RandomInt16 generates a random integer between min and max
func RandomInt16() int16 {
	return int16(rand.Int())
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func ConvertStringInt16(s string) (int16, error) {
	i, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return -1, err
	}
	return int16(i), nil
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
