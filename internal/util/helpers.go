package util

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

var (
	Integration = false
)

// TranslatesReturn avoiding cycle between mock/translates and translates/translates
type TranslatesReturn struct {
	PhraseId int64
	Text     string
}

// HashPassword returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// PathExists returns whether the given file or directory exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ConvertStringInt16(s string) (int16, error) {
	i, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return -1, err
	}
	return int16(i), nil
}

func SafeCastToInt16(value int) (int16, error) {
	if value < math.MinInt16 || value > math.MaxInt16 {
		return 0, fmt.Errorf("value %d is out of range for int16", value)
	}
	return int16(value), nil
}
