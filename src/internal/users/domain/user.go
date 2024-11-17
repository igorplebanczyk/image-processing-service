package domain

import (
	"fmt"
	"github.com/google/uuid"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

type Role string

func (r Role) String() string {
	return string(r)
}

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(username, email, password string) *User {
	return &User{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  password,
		Role:      UserRole,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if len(username) < 3 || len(username) > 32 {
		return fmt.Errorf("username must be between 3 and 32 characters")
	}

	if strings.Contains(username, " ") {
		return fmt.Errorf("username cannot contain spaces")
	}

	return nil
}

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if len(password) < 8 || len(password) > 32 {
		return fmt.Errorf("password must be between 8 and 32 characters")
	}

	if strings.Contains(password, " ") {
		return fmt.Errorf("password cannot contain spaces")
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

func DetermineUserDetailsToUpdate(existingUser *User, newUsername, newEmail string) (string, string, error) {
	if newUsername == "" && newEmail == "" {
		return "", "", fmt.Errorf("username or email must be provided")
	}

	if newUsername == "" {
		newUsername = existingUser.Username
	}

	if newEmail == "" {
		newEmail = existingUser.Email
	}

	return newUsername, newEmail, nil
}
