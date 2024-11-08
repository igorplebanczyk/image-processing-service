package images

import "github.com/google/uuid"

type Repository interface {
	CreateImage(userID uuid.UUID, name string) (*Image, error)
	GetImagesByUserID(userID uuid.UUID) ([]*Image, error)
	UpdateImage(userID uuid.UUID) error
	DeleteImage(userID uuid.UUID, name string) error
}
