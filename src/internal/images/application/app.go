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
	dbRepo                domain.ImageDBRepository
	storageRepo           domain.ImageStorageRepository
	cacheRepo             domain.ImageCacheRepository
	transformationService *transformations.Service
	cacheExpiry           time.Duration
}

func NewService(
	dbRepo domain.ImageDBRepository,
	storageRepo domain.ImageStorageRepository,
	cacheRepo domain.ImageCacheRepository,
	cacheExpiry time.Duration,
) *ImageService {
	return &ImageService{
		dbRepo:                dbRepo,
		storageRepo:           storageRepo,
		cacheRepo:             cacheRepo,
		transformationService: transformations.NewService(),
		cacheExpiry:           cacheExpiry,
	}
}

func (s *ImageService) Upload(userID uuid.UUID, name, description string, bytes []byte) error {
	err := domain.ValidateName(name)
	if err != nil {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid image name: %v", err))
	}

	err = domain.ValidateDescription(description)
	if err != nil {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid image description: %v", err))
	}

	err = domain.ValidateImage(bytes)
	if err != nil {
		return commonerrors.NewInvalidInput("invalid image data")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageMetadata, err := s.dbRepo.CreateImageMetadata(ctx, userID, name, description)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error creating image in database: %v", err))
	}

	fullImageObjectName := domain.CreateFullImageObjectName(imageMetadata.ID)
	err = s.storageRepo.UploadImage(ctx, fullImageObjectName, bytes)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error uploading image to storage: %v", err))
	}
	err = s.cacheRepo.CacheImage(ctx, fullImageObjectName, bytes, s.cacheExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
	}

	previewBytes, err := s.transformationService.CreatePreview(bytes)
	previewImageObjectName := domain.CreatePreviewImageObjectName(imageMetadata.ID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error creating preview image: %v", err))
	}
	err = s.storageRepo.UploadImage(ctx, previewImageObjectName, previewBytes)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error uploading preview image to storage: %v", err))
	}
	err = s.cacheRepo.CacheImage(ctx, previewImageObjectName, previewBytes, s.cacheExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error caching preview image: %v", err))
	}

	return nil
}

func (s *ImageService) Get(userID uuid.UUID, name string) (*domain.ImageMetadata, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageMetadata, err := s.dbRepo.GetImageMetadataByUserIDAndName(ctx, userID, name)
	if err != nil {
		return nil, nil, commonerrors.NewInternal(fmt.Sprintf("error reading image metadata from database: %v", err))
	}

	fullImageObjectName := domain.CreateFullImageObjectName(imageMetadata.ID)
	imageBytes, err := s.cacheRepo.GetImage(ctx, fullImageObjectName)
	if err != nil {
		return nil, nil, commonerrors.NewInternal(fmt.Sprintf("error reading image from cache: %v", err))
	}
	if imageBytes == nil {
		imageBytes, err = s.storageRepo.DownloadImage(ctx, fullImageObjectName)
		if err != nil {
			return nil, nil, commonerrors.NewInternal(fmt.Sprintf("error downloading image from storage: %v", err))
		}
		err = s.cacheRepo.CacheImage(ctx, fullImageObjectName, imageBytes, s.cacheExpiry)
		if err != nil {
			return nil, nil, commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
		}
	}

	return imageMetadata, imageBytes, nil
}

func (s *ImageService) GetAll(userID uuid.UUID, page, limit int) ([]*domain.ImageMetadata, [][]byte, int, error) {
	if page < 1 || limit < 1 || limit > 25 {
		return nil, nil, -1, commonerrors.NewInvalidInput("invalid page or limit")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imagesMetadata, totalCount, err := s.dbRepo.GetImagesMetadataByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, nil, -1, commonerrors.NewInternal(fmt.Sprintf("error reading images metadata from database: %v", err))
	}

	imagesBytes := make([][]byte, len(imagesMetadata))
	for i, imageMetadata := range imagesMetadata {
		previewImageObjectName := domain.CreatePreviewImageObjectName(imageMetadata.ID)
		imageBytes, err := s.cacheRepo.GetImage(ctx, previewImageObjectName)
		if err != nil {
			return nil, nil, -1, commonerrors.NewInternal(fmt.Sprintf("error reading image from cache: %v", err))
		}
		if imageBytes == nil {
			imageBytes, err = s.storageRepo.DownloadImage(ctx, previewImageObjectName)
			if err != nil {
				return nil, nil, -1, commonerrors.NewInternal(fmt.Sprintf("error downloading image from storage: %v", err))
			}
			err = s.cacheRepo.CacheImage(ctx, previewImageObjectName, imageBytes, s.cacheExpiry)
			if err != nil {
				return nil, nil, -1, commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
			}
		}
		imagesBytes[i] = imageBytes
	}

	return imagesMetadata, imagesBytes, totalCount, nil
}

func (s *ImageService) UpdateDetails(userID uuid.UUID, oldName, newName, newDescription string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageMetadata, err := s.dbRepo.GetImageMetadataByUserIDAndName(ctx, userID, oldName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error reading image metadata from database: %v", err))
	}

	newName, newDescription, err = domain.DetermineImageMetadataToUpdate(imageMetadata, newName, newDescription)
	if err != nil {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid image metadata: %v", err))
	}

	err = domain.ValidateName(newName)
	if err != nil {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid image name: %v", err))
	}

	err = domain.ValidateDescription(newDescription)
	if err != nil {
		return commonerrors.NewInvalidInput(fmt.Sprintf("invalid image description: %v", err))
	}

	err = s.dbRepo.UpdateImageMetadataDetails(ctx, imageMetadata.ID, newName, newDescription)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error updating image metadata in database: %v", err))
	}

	return nil
}

func (s *ImageService) Transform(
	userID uuid.UUID,
	name string,
	transformations []domain.Transformation,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageMetadata, err := s.dbRepo.GetImageMetadataByUserIDAndName(ctx, userID, name)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error reading image metadata from database: %v", err))
	}

	fullImageObjectName := domain.CreateFullImageObjectName(imageMetadata.ID)
	imageBytes, err := s.cacheRepo.GetImage(ctx, fullImageObjectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error reading image from cache: %v", err))
	}
	if imageBytes == nil {
		imageBytes, err = s.storageRepo.DownloadImage(ctx, fullImageObjectName)
		if err != nil {
			return commonerrors.NewInternal(fmt.Sprintf("error downloading image from storage: %v", err))
		}
		err = s.cacheRepo.CacheImage(ctx, fullImageObjectName, imageBytes, s.cacheExpiry)
		if err != nil {
			return commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
		}
	}

	transformedBytes, err := s.transformationService.Apply(imageBytes, transformations)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error applying transformations: %v", err))
	}

	err = s.storageRepo.UploadImage(ctx, fullImageObjectName, transformedBytes)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error uploading image to storage: %v", err))
	}
	err = s.cacheRepo.CacheImage(ctx, fullImageObjectName, transformedBytes, s.cacheExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
	}

	previewBytes, err := s.transformationService.CreatePreview(transformedBytes)
	previewImageObjectName := domain.CreatePreviewImageObjectName(imageMetadata.ID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error creating preview image: %v", err))
	}
	err = s.storageRepo.UploadImage(ctx, previewImageObjectName, previewBytes)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error uploading preview image to storage: %v", err))
	}
	err = s.cacheRepo.CacheImage(ctx, previewImageObjectName, previewBytes, s.cacheExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error caching preview image: %v", err))
	}

	err = s.dbRepo.UpdateImageMetadataUpdatedAt(ctx, imageMetadata.ID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error updating image metadata in database: %v", err))
	}

	return nil
}

func (s *ImageService) Delete(userID uuid.UUID, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageMetadata, err := s.dbRepo.GetImageMetadataByUserIDAndName(ctx, userID, name)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error reading image metadata from database: %v", err))
	}

	err = s.dbRepo.DeleteImageMetadata(ctx, imageMetadata.ID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting image metadata from database: %v", err))
	}

	fullImageObjectName := domain.CreateFullImageObjectName(imageMetadata.ID)
	err = s.storageRepo.DeleteImage(ctx, fullImageObjectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting image from storage: %v", err))
	}

	err = s.cacheRepo.DeleteImage(ctx, fullImageObjectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting image from cache: %v", err))
	}

	previewImageObjectName := domain.CreatePreviewImageObjectName(imageMetadata.ID)
	err = s.storageRepo.DeleteImage(ctx, previewImageObjectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting preview image from storage: %v", err))
	}

	err = s.cacheRepo.DeleteImage(ctx, previewImageObjectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting preview image from cache: %v", err))
	}

	return nil
}
