package domain

import (
	"context"
	"github.com/google/uuid"
)

type ImageDBRepository interface {
	CreateImageMetadata(
		ctx context.Context,
		userID uuid.UUID,
		name,
		description string,
	) (*ImageMetadata, error)
	GetImageMetadataByUserIDAndName(ctx context.Context, userID uuid.UUID, name string) (*ImageMetadata, error)
	GetImagesMetadataByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*ImageMetadata, int, error)
	GetAllImagesMetadata(ctx context.Context, page, limit int) ([]*ImageMetadata, int, error)
	UpdateImageMetadataDetails(
		ctx context.Context,
		id uuid.UUID,
		newName,
		newDescription string,
	) error
	UpdateImageMetadataUpdatedAt(ctx context.Context, id uuid.UUID) error
	DeleteImageMetadata(ctx context.Context, id uuid.UUID) error
}
