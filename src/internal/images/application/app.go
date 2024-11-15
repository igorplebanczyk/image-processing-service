package application

import (
	"context"
	"errors"
	"github.com/google/uuid"
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
		return nil, errors.Join(domain.ErrValidationFailed, err)
	}

	err = domain.ValidateRawImage(imageBytes)
	if err != nil {
		return nil, errors.Join(domain.ErrValidationFailed, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := s.repo.CreateImage(ctx, userID, imageName)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	objectName := domain.CreateObjectName(userID, imageName)
	err = s.storage.Upload(ctx, objectName, imageBytes)
	if err != nil {
		cancel()
		return nil, errors.Join(domain.ErrInternal, err)
	}

	err = s.cache.Set(ctx, objectName, imageBytes, time.Hour)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	return image, nil
}

func (s *ImageService) ListUserImages(userID uuid.UUID, page, limit *int) ([]*domain.Image, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	images, total, err := s.repo.GetImagesByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, errors.Join(domain.ErrInternal, err)
	}

	return images, total, nil
}

func (s *ImageService) GetImageData(userID uuid.UUID, imageName string) (*domain.Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := s.repo.GetImageByUserIDandName(ctx, userID, imageName)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
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
			return nil, errors.Join(domain.ErrInternal, err)
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
		return errors.Join(domain.ErrInternal, err)
	}

	err = s.repo.DeleteImage(ctx, image.ID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	err = s.storage.Delete(ctx, objectName)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	err = s.cache.Delete(ctx, objectName)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
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

	for _, transformation := range transformations {
		imageBytes, err = s.transformationService.Apply(imageBytes, transformation)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidRequest) {
				return errors.Join(domain.ErrInvalidRequest, err)
			}
			return err
		}
	}

	objectName := domain.CreateObjectName(userID, imageName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := s.repo.GetImageByUserIDandName(ctx, userID, imageName)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	err = s.repo.UpdateImage(ctx, image.ID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	err = s.storage.Upload(ctx, objectName, imageBytes)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	err = s.cache.Set(ctx, objectName, imageBytes, time.Hour)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}
