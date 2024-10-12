package utils

import (
	"errors"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func ValidateEmail(email string) []error {
	err := validator.New().Struct(struct {
		Email string `validate:"required,email"`
	}{email})
	if err != nil {
		return []error{errors.New("Email must match pattern")}
	}
	return []error{}
}

func ValidatePassword(password string) []error {

	var hasUpper, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	errorMessages := []error{}
	if len(password) < 8 {
		errorMessages = append(errorMessages, errors.New("Password must be at least 8 characters long"))
	}
	if !hasUpper {
		errorMessages = append(errorMessages, errors.New("Password must contain at least one uppercase letter"))
	}
	if !hasDigit {
		errorMessages = append(errorMessages, errors.New("Password must contain at least one digit"))
	}
	if !hasSpecial {
		errorMessages = append(errorMessages, errors.New("Password must contain at least one special character"))
	}

	return errorMessages
}
