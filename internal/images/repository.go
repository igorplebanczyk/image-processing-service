package images

import "github.com/google/uuid"

type Repository interface {
	CreateImage(userID uuid.UUID, name string) (*Image, error)
	DeleteImage(userID uuid.UUID, name string) error
	UpdateImage(userID uuid.UUID) error
	GetImagesByUserID(userID uuid.UUID) ([]*Image, error)
}
