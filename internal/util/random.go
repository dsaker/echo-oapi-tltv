package util

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
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

// RandomName generates a random user name
func RandomName(n int) string {
	return RandomString(n)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
