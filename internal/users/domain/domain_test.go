package domain

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		username string
		wantErr  bool
	}{
		{"validuser", false},
		{"", true},
		{"ab", true},
		{"averyverylongusernamethatiswaytoolong", true},
	}

	for _, tt := range tests {
		err := ValidateUsername(tt.username)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateUsername(%v) error = %v, wantErr %v", tt.username, err, tt.wantErr)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email   string
		wantErr bool
	}{
		{"test@example.com", false},
		{"invalid-email", true},
		{"@missingusername.com", true},
	}

	for _, tt := range tests {
		err := ValidateEmail(tt.email)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateEmail(%v) error = %v, wantErr %v", tt.email, err, tt.wantErr)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		wantErr  bool
	}{
		{"Valid@123", false},
		{"short", true},
		{"loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong", true},
		{"noSpecialCharacter123", true},
		{"NOLOWERCASE@123", true},
		{"nouppercase@123", true},
		{"noNumber@", true},
	}

	for _, tt := range tests {
		err := ValidatePassword(tt.password)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePassword(%v) error = %v, wantErr %v", tt.password, err, tt.wantErr)
		}
	}
}

func TestHashPassword(t *testing.T) {
	password := "SecureP@ss123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		t.Errorf("expected hashed password to match original password, got error = %v", err)
	}
	if hashedPassword == password {
		t.Errorf("expected hashed password to be different from original password")
	}
}
