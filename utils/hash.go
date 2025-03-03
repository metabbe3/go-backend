package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher defines an interface for password hashing
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword, plainPassword string) error
}

// BcryptHasher implements PasswordHasher using bcrypt
type BcryptHasher struct{}

// HashPassword hashes a password using bcrypt
func (BcryptHasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswords checks if the provided password matches the hashed password
func (BcryptHasher) ComparePasswords(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
