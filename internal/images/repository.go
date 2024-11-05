package images

import "github.com/google/uuid"

type Repository interface {
	CreateImage(userID uuid.UUID, name string) (*Image, error)
	GetImageByID(id uuid.UUID) (*Image, error)
	GetImagesByUserID(userID uuid.UUID) ([]*Image, error)
}
