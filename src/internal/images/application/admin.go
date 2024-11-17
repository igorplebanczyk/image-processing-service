package application

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"image-processing-service/src/internal/images/domain"
	"time"
)

func (s *ImageService) AdminListAllImages(page, limit *int) ([]*domain.Image, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	images, total, err := s.repo.GetAllImages(ctx, page, limit)
	if err != nil {
		return nil, 0, errors.Join(domain.ErrInternal, err)
	}

	return images, total, nil
}

func (s *ImageService) AdminDeleteImage(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.repo.DeleteImage(ctx, id)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}
