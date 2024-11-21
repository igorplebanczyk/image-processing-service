package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	commonerrors "image-processing-service/src/internal/common/errors"
	"image-processing-service/src/internal/images/domain"
	"time"
)

func (s *ImageService) AdminListAllImages(page, limit *int) ([]*domain.Image, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	images, total, err := s.repo.GetAllImages(ctx, page, limit)
	if err != nil {
		return nil, 0, commonerrors.NewInternal(fmt.Sprintf("error fetching images from database: %v", err))
	}

	return images, total, nil
}

func (s *ImageService) AdminDeleteImage(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.repo.DeleteImage(ctx, id)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting image from database: %v", err))
	}

	return nil
}
