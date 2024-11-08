package users

import (
	"fmt"
	"net/mail"
	"regexp"
)

func validateUsername(r UserRepository, username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	user, _ := r.GetUserByUsername(username)
	if user != nil {
		return fmt.Errorf("users with username %s already exists", username)
	}

	return nil
}

func validateEmail(r UserRepository, email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email address")
	}

	user, _ := r.GetUserByEmail(email)
	if user != nil {
		return fmt.Errorf("users with email %s already exists", email)
	}

	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if len(password) < 8 || len(password) > 32 {
		return fmt.Errorf("password must be between 8 and 32 characters")
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
	hasDigit := regexp.MustCompile(`\d`).MatchString
	hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString

	if !(hasLower(password) && hasUpper(password) && hasDigit(password) && hasSpecial(password)) {
		return fmt.Errorf("password must contain at least one lowercase letter, one uppercase letter, one digit, and one special character")
	}

	return nil
}
