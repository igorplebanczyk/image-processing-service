package domain

import (
	"github.com/google/uuid"
	"strings"
	"testing"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"ValidName", "John", false},
		{"TooShortName", "Jo", true},
		{"TooLongName", string(make([]byte, 129)), true},
		{"MinimumLengthName", "abc", false},
		{"MaximumLengthName", string(make([]byte, 128)), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateName(tc.input)
			if tc.expectError && err == nil {
				t.Fatalf("expected an error, got none")
			}
			if !tc.expectError && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidateRawImage(t *testing.T) {
	tests := []struct {
		name        string
		imageBytes  []byte
		expectError bool
	}{
		{"ValidImageSize", make([]byte, MaxImageSize-1), false},
		{"ExactMaxSizeImage", make([]byte, MaxImageSize), false},
		{"ExceedsMaxSize", make([]byte, MaxImageSize+1), true},
		{"EmptyImage", make([]byte, 0), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateRawImage(tc.imageBytes)
			if tc.expectError && err == nil {
				t.Fatalf("expected an error, got none")
			}
			if !tc.expectError && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestCreateObjectName(t *testing.T) {
	userID := uuid.New()
	imageName := "example.jpg"
	objectName := CreateObjectName(userID, imageName)

	if objectName == "" {
		t.Fatalf("expected a non-empty object name, got empty string")
	}

	if !strings.Contains(objectName, userID.String()) {
		t.Errorf("expected object name to contain user ID, got %s", objectName)
	}

	if !strings.Contains(objectName, imageName) {
		t.Errorf("expected object name to contain image name, got %s", objectName)
	}
}
