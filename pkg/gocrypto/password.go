package gocrypto

import (
	"golang.org/x/crypto/bcrypt"
)

// HashAndSaltPassword hash password with salt
func HashAndSaltPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

// VerifyPassword verify password and ciphertext match
func VerifyPassword(password string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
