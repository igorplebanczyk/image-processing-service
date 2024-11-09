package images

import "github.com/google/uuid"

type Repository interface {
	CreateImage(userID uuid.UUID, name string) (*Image, error)
	GetImagesByUserID(userID uuid.UUID) ([]*Image, error)
	GetImageByUserIDandName(userID uuid.UUID, name string) (*Image, error)
	UpdateImage(id uuid.UUID) error
	DeleteImage(id uuid.UUID) error
}
