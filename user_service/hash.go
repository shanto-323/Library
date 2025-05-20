package userservice

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CreateNewHashPassword(password string) (string, error) {
	h_pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(h_pass), nil
}

func CompareWithHash(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("password not matched %s", err)
	}
	return nil
}
