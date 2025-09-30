package utils

import (
	"errors"
	"regexp"
)

// ValidatePasswordStrength checks if the password meets strength requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	var (
		upper   = regexp.MustCompile(`[A-Z]`)
		lower   = regexp.MustCompile(`[a-z]`)
		digit   = regexp.MustCompile(`[0-9]`)
		special = regexp.MustCompile(`[!@#\$%\^&\*\-_]`)
	)
	if !upper.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !lower.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !digit.MatchString(password) {
		return errors.New("password must contain at least one digit")
	}
	if !special.MatchString(password) {
		return errors.New("password must contain at least one special character (!@#$%^&*-_)")
	}
	return nil
}
