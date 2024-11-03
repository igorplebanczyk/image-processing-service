package database

import (
	"fmt"
	"regexp"
)

func validate(r *UserRepository, username, email, password string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if r.GetUserByUsername(username) != nil {
		return fmt.Errorf("user with username %s already exists", username)
	}

	if r.GetUserByEmail(email) != nil {
		return fmt.Errorf("user with email %s already exists", email)
	}

	match, err := regexp.MatchString(passwordRegex, password)
	if err != nil {
		return fmt.Errorf("error validating password: %w", err)
	}
	if !match {
		return fmt.Errorf("password must contain at least one lowercase letter, one uppercase letter, one digit, one special character, and be between 8 and 32 characters long")
	}

	return nil
}
