package utils

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using the validator tags
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}

// IsValidEmail checks if the email format is valid
func IsValidEmail(email string) bool {
	// RFC 5322 compliant email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidUsername checks if the username is valid (alphanumeric + underscore, 3-30 chars)
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// IsStrongPassword checks if the password meets strength requirements
// Requirements: min 8 chars, at least 1 uppercase, 1 lowercase, 1 number
func IsStrongPassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasLower {
		return false, "Password must contain at least one lowercase letter"
	}
	if !hasNumber {
		return false, "Password must contain at least one number"
	}

	return true, ""
}

// ValidatePassword validates password strength and returns an error if invalid
func ValidatePassword(password string) error {
	valid, message := IsStrongPassword(password)
	if !valid {
		return fmt.Errorf(message)
	}
	return nil
}
