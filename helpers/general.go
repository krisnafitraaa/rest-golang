package helpers

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return string(bytes), err
	}

	return string(bytes), nil
}

func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return false
	}

	return true
}

func MakeRandomString(length int) string {
	str := make([]rune, length)

	for i := range str {
		str[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(str)
}
