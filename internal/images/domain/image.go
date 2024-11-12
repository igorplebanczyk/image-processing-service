package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

const MaxImageSize = 10 * 1024 * 1024

type Image struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewImage(userID uuid.UUID, name string) *Image {
	return &Image{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func ValidateName(name string) error {
	if len(name) < 3 || len(name) > 128 {
		return fmt.Errorf("name must be between 3 and 128 characters")
	}

	return nil
}

func ValidateRawImage(imageBytes []byte) error {
	if len(imageBytes) > MaxImageSize {
		return fmt.Errorf("image size cannot exceed %d bytes", MaxImageSize)
	}

	return nil
}

func CreateObjectName(userID uuid.UUID, imageName string) string {
	return fmt.Sprintf("%s-%s", userID, imageName)
}
