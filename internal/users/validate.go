package users

import (
	"fmt"
	"image-processing-service/internal/services/database"
	"regexp"
)

func validate(r *database.UserRepository, username, email, password string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	user, _ := r.GetUserByValue("username", username)
	if user != nil {
		return fmt.Errorf("users with username %s already exists", username)
	}

	user, _ = r.GetUserByValue("email", email)
	if user != nil {
		return fmt.Errorf("users with email %s already exists", email)
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
