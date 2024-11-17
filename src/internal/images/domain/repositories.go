package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type ImageRepository interface {
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
	GetImagesByUserID(
		ctx context.Context,
		userID uuid.UUID,
		page,
		limit *int,
	) ([]*Image, int, error)
	UpdateImage(ctx context.Context, id uuid.UUID) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
	GetAllImages(ctx context.Context, page, limit *int) ([]*Image, int, error)
}

type ImageStorageRepository interface {
	Upload(ctx context.Context, objectName string, data []byte) error
	Download(ctx context.Context, objectName string) ([]byte, error)
	Delete(ctx context.Context, objectName string) error
}

type ImageCacheRepository interface {
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}
