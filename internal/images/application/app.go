package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/images/application/transformations"
	"image-processing-service/internal/images/domain"
	"log/slog"
	"net/http"
	"time"
)

type ImageService struct {
	repo                  domain.ImageRepository
	storage               domain.StorageService
	cache                 domain.CacheService
	transformationService *transformations.Service
}

func NewService(
	repo domain.ImageRepository,
	storage domain.StorageService,
	cache domain.CacheService,
	transformationWorkerCount int,
	transformationQueueSize int,
) *ImageService {
	return &ImageService{
		repo:                  repo,
		storage:               storage,
		cache:                 cache,
		transformationService: transformations.New(transformationWorkerCount, transformationQueueSize),
	}
}

func (s *ImageService) UploadImage(userID uuid.UUID, imageName string, imageBytes []byte) (*domain.Image, error) {
	err := domain.ValidateName(imageName)
	if err != nil {
		return nil, fmt.Errorf("invalid image name: %w", err)
	}

	err = domain.ValidateRawImage(imageBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid image: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := s.repo.CreateImage(ctx, userID, imageName)
	if err != nil {
		return nil, fmt.Errorf("error creating image: %w", err)
	}

	objectName := domain.CreateObjectName(userID, imageName)
	err = s.storage.Upload(ctx, objectName, imageBytes)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("error uploading image: %w", err)
	}

	err = s.cache.Set(objectName, imageBytes, time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error caching image: %w", err)
	}

	return image, nil
}

func (s *ImageService) ListUserImages(userID uuid.UUID, page, limit *int) ([]*domain.Image, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	images, total, err := s.repo.GetImagesByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching images: %w", err)
	}

	return images, total, nil
}

func (s *ImageService) GetImageData(userID uuid.UUID, imageName string) (*domain.Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := s.repo.GetImageByUserIDandName(ctx, userID, imageName)
	if err != nil {
		return nil, fmt.Errorf("error getting image: %w", err)
	}

	return image, nil
}

func (s *ImageService) DownloadImage(userID uuid.UUID, imageName string) ([]byte, error) {
	objectName := domain.CreateObjectName(userID, imageName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageBytes, err := s.cache.Get(objectName)
	if err != nil {
		imageBytes, err = s.storage.Download(ctx, objectName)
		if err != nil {
			return nil, fmt.Errorf("error downloading image: %w", err)
		}
	}
	slog.Info(fmt.Sprintf("format at downloadImage: %v", http.DetectContentType(imageBytes)))

	return imageBytes, nil
}

func (s *ImageService) DeleteImage(userID uuid.UUID, imageName string) error {
	objectName := domain.CreateObjectName(userID, imageName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := s.repo.GetImageByUserIDandName(ctx, userID, imageName)
	if err != nil {
		return fmt.Errorf("error getting image ID: %w", err)
	}

	err = s.repo.DeleteImage(ctx, image.ID)
	if err != nil {
		return fmt.Errorf("error deleting image: %w", err)
	}

	err = s.storage.Delete(ctx, objectName)
	if err != nil {
		return fmt.Errorf("error deleting image: %w", err)
	}

	err = s.cache.Delete(objectName)
	if err != nil {
		return fmt.Errorf("error deleting image: %w", err)
	}

	return nil
}

func (s *ImageService) ApplyTransformations(
	userID uuid.UUID,
	imageName string,
	transformations []domain.Transformation,
) error {
	imageBytes, err := s.DownloadImage(userID, imageName)
	if err != nil {
		return fmt.Errorf("error downloading image: %w", err)
	}

	for _, transformation := range transformations {
		imageBytes, err = s.transformationService.Apply(imageBytes, transformation)
		if err != nil {
			return fmt.Errorf("error applying transformation: %w", err)
		}
	}

	objectName := domain.CreateObjectName(userID, imageName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := s.repo.GetImageByUserIDandName(ctx, userID, imageName)
	if err != nil {
		return fmt.Errorf("error getting image: %w", err)
	}

	err = s.repo.UpdateImage(ctx, image.ID)
	if err != nil {
		return fmt.Errorf("error updating image: %w", err)
	}

	err = s.storage.Upload(ctx, objectName, imageBytes)
	if err != nil {
		return fmt.Errorf("error uploading transformed image: %w", err)
	}

	err = s.cache.Set(objectName, imageBytes, time.Hour)
	if err != nil {
		return fmt.Errorf("error caching transformed image: %w", err)
	}

	return nil
}
