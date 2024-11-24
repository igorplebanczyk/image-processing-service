package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	commonerrors "image-processing-service/src/internal/common/errors"
	"image-processing-service/src/internal/images/application/transformations"
	"image-processing-service/src/internal/images/domain"
	"time"
)

type ImageService struct {
	repo                  domain.ImageRepository
	storage               domain.ImageStorageRepository
	cache                 domain.ImageCacheRepository
	transformationService *transformations.Service
}

func NewService(
	repo domain.ImageRepository,
	storage domain.ImageStorageRepository,
	cache domain.ImageCacheRepository,
) *ImageService {
	return &ImageService{
		repo:                  repo,
		storage:               storage,
		cache:                 cache,
		transformationService: transformations.NewService(),
	}
}

func (s *ImageService) UploadImage(userID uuid.UUID, imageName string, imageBytes []byte) (*domain.Image, error) {
	err := domain.ValidateName(imageName)
	if err != nil {
		return nil, commonerrors.NewInvalidInput(fmt.Sprintf("invalid image name: %v", err))
	}

	err = domain.ValidateRawImage(imageBytes)
	if err != nil {
		return nil, commonerrors.NewInvalidInput("invalid image data")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := s.repo.CreateImage(ctx, userID, imageName)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("error creating image in database: %v", err))
	}

	objectName := domain.CreateObjectName(userID, imageName)
	err = s.storage.Upload(ctx, objectName, imageBytes)
	if err != nil {
		cancel()
		return nil, commonerrors.NewInternal(fmt.Sprintf("error uploading image to storage: %v", err))
	}

	err = s.cache.Set(ctx, objectName, imageBytes, time.Hour)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
	}

	return image, nil
}

func (s *ImageService) ListUserImages(userID uuid.UUID, page, limit *int) ([]*domain.Image, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	images, total, err := s.repo.GetImagesByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, commonerrors.NewInternal(fmt.Sprintf("error fetching images from database: %v", err))
	}

	return images, total, nil
}

func (s *ImageService) GetImageData(userID uuid.UUID, imageName string) (*domain.Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := s.repo.GetImageByUserIDandName(ctx, userID, imageName)
	if err != nil {
		return nil, commonerrors.NewInternal(fmt.Sprintf("error reading image from database: %v", err))
	}

	return image, nil
}

func (s *ImageService) DownloadImage(userID uuid.UUID, imageName string) ([]byte, error) {
	objectName := domain.CreateObjectName(userID, imageName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageBytes, err := s.cache.Get(ctx, objectName)
	if err != nil {
		imageBytes, err = s.storage.Download(ctx, objectName)
		if err != nil {
			return nil, commonerrors.NewInternal(fmt.Sprintf("error downloading image from storage: %v", err))
		}
	}

	return imageBytes, nil
}

func (s *ImageService) DeleteImage(userID uuid.UUID, imageName string) error {
	objectName := domain.CreateObjectName(userID, imageName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := s.repo.GetImageByUserIDandName(ctx, userID, imageName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error reading image from database: %v", err))
	}

	err = s.repo.DeleteImage(ctx, image.ID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting image from database: %v", err))
	}

	err = s.storage.Delete(ctx, objectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting image from storage: %v", err))
	}

	err = s.cache.Delete(ctx, objectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting image from cache: %v", err))
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
		return err
	}

	imageBytes, err = s.transformationService.Apply(imageBytes, transformations)
	if err != nil {
		return err
	}

	objectName := domain.CreateObjectName(userID, imageName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := s.repo.GetImageByUserIDandName(ctx, userID, imageName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error reading image from database: %v", err))
	}

	err = s.repo.UpdateImage(ctx, image.ID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error updating image in database: %v", err))
	}

	err = s.storage.Upload(ctx, objectName, imageBytes)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error uploading image to storage: %v", err))
	}

	err = s.cache.Set(ctx, objectName, imageBytes, time.Hour)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
	}

	return nil
}
