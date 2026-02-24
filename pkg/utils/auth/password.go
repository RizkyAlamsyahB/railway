package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// DefaultCost is the bcrypt cost factor used for hashing passwords.
const DefaultCost = bcrypt.DefaultCost

// HashPassword returns a bcrypt hash of the given password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// CheckPassword compares a plaintext password against a bcrypt hash.
// Returns nil if they match, or an error if they don't.
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
