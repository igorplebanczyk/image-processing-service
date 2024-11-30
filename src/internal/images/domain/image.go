package domain

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

const MaxImageSize = 10 * 1024 * 1024

type ImageMetadata struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewImageMetadata(userID uuid.UUID, name, description string) *ImageMetadata {
	return &ImageMetadata{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func ValidateName(name string) error {
	if len(name) < 3 || len(name) > 128 {
		return fmt.Errorf("name must be between 3 and 128 characters")
	}

	if strings.Contains(name, " ") {
		return fmt.Errorf("name cannot contain spaces")
	}

	return nil
}

func ValidateDescription(description string) error {
	if len(description) > 1024 {
		return fmt.Errorf("description cannot exceed 1024 characters")
	}

	return nil
}

func ValidateImage(bytes []byte) error {
	if len(bytes) > MaxImageSize {
		return fmt.Errorf("image size cannot exceed %d bytes", MaxImageSize)
	}

	return nil
}

func DetermineImageMetadataToUpdate(existingImageMetadata *ImageMetadata, newName, newDescription string) (string, string, error) {
	if newName == "" && newDescription == "" {
		return "", "", fmt.Errorf("no fields to update")
	}

	if newName == "" {
		newName = existingImageMetadata.Name
	}

	if newDescription == "" {
		newDescription = existingImageMetadata.Description
	}

	return newName, newDescription, nil
}

func CreateFullImageObjectName(id uuid.UUID) string {
	return fmt.Sprintf("full-%s", id)
}

func CreatePreviewImageObjectName(id uuid.UUID) string {
	return fmt.Sprintf("prev-%s", id)
}
