package utils

import (
	"errors"
	"regexp"
	"strings"
)

// IsValidEmail checks if an email is valid
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// IsValidPassword checks if a password is strong (min 8 chars, 1 uppercase, 1 number)
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasDigit := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasDigit
}

// TrimString trims leading & trailing spaces from a string
func TrimString(input string) string {
	return strings.TrimSpace(input)
}

// ValidateUserInput checks if email and password are valid
func ValidateUserInput(email, password string) error {
	if !IsValidEmail(email) {
		return errors.New("invalid email format")
	}

	if !IsValidPassword(password) {
		return errors.New("password must be at least 8 characters, include 1 uppercase letter and 1 number")
	}

	return nil
}
