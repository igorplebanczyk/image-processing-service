package images

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreateImage(
		ctx context.Context,
		userID uuid.UUID,
		name string,
	) (*Image, error)
	GetImageByUserIDandName(
		ctx context.Context,
		userID uuid.UUID,
		name string,
	) (*Image, error)
	GetImagesByUserID(ctx context.Context, userID uuid.UUID) ([]*Image, error)
	UpdateImage(ctx context.Context, id uuid.UUID) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
}
