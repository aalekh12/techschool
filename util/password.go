package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to hash password %s", err)
	}
	return string(hashpass), nil
}

func ComparePassword(password string, HashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(HashPassword), []byte(password))
}
