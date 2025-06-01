package crypt

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	const op = "crypt.HashPassword"
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passHash), fmt.Errorf("%s: %w", op, err)
}

func CheckPasswordHash(password string, passHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password))
	return err == nil
}
