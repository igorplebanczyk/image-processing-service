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

type ImagesService struct {
	imagesDBRepo           domain.ImagesDBRepository
	imagesStorageRepo      domain.ImagesStorageRepository
	imagesCacheRepo        domain.ImagesCacheRepository
	transformationsService *transformations.Service
	cacheExpiry            time.Duration
}

func NewService(
	imagesDBRepo domain.ImagesDBRepository,
	imagesStorageRepo domain.ImagesStorageRepository,
	imagesCacheRepo domain.ImagesCacheRepository,
	cacheExpiry time.Duration,
) *ImagesService {
	return &ImagesService{
		imagesDBRepo:           imagesDBRepo,
		imagesStorageRepo:      imagesStorageRepo,
		imagesCacheRepo:        imagesCacheRepo,
		transformationsService: transformations.NewService(),
		cacheExpiry:            cacheExpiry,
	}
}

func (s *ImagesService) Upload(userID uuid.UUID, name, description string, bytes []byte) error {
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

	imageMetadata, err := s.imagesDBRepo.CreateImageMetadata(ctx, userID, name, description)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error creating image in database: %v", err))
	}

	fullImageObjectName := domain.CreateFullImageObjectName(imageMetadata.ID)
	err = s.imagesStorageRepo.UploadImage(ctx, fullImageObjectName, bytes)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error uploading image to storage: %v", err))
	}
	err = s.imagesCacheRepo.CacheImage(ctx, fullImageObjectName, bytes, s.cacheExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
	}

	previewBytes, err := s.transformationsService.CreatePreview(bytes)
	previewImageObjectName := domain.CreatePreviewImageObjectName(imageMetadata.ID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error creating preview image: %v", err))
	}
	err = s.imagesStorageRepo.UploadImage(ctx, previewImageObjectName, previewBytes)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error uploading preview image to storage: %v", err))
	}
	err = s.imagesCacheRepo.CacheImage(ctx, previewImageObjectName, previewBytes, s.cacheExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error caching preview image: %v", err))
	}

	return nil
}

func (s *ImagesService) Get(userID uuid.UUID, name string) (*domain.ImageMetadata, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageMetadata, err := s.imagesDBRepo.GetImageMetadataByUserIDAndName(ctx, userID, name)
	if err != nil {
		return nil, nil, commonerrors.NewInternal(fmt.Sprintf("error reading image metadata from database: %v", err))
	}

	fullImageObjectName := domain.CreateFullImageObjectName(imageMetadata.ID)
	imageBytes, err := s.imagesCacheRepo.GetImage(ctx, fullImageObjectName)
	if err != nil {
		return nil, nil, commonerrors.NewInternal(fmt.Sprintf("error reading image from cache: %v", err))
	}
	if imageBytes == nil {
		imageBytes, err = s.imagesStorageRepo.DownloadImage(ctx, fullImageObjectName)
		if err != nil {
			return nil, nil, commonerrors.NewInternal(fmt.Sprintf("error downloading image from storage: %v", err))
		}
		err = s.imagesCacheRepo.CacheImage(ctx, fullImageObjectName, imageBytes, s.cacheExpiry)
		if err != nil {
			return nil, nil, commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
		}
	}

	return imageMetadata, imageBytes, nil
}

func (s *ImagesService) GetAll(userID uuid.UUID, page, limit int) ([]*domain.ImageMetadata, [][]byte, int, error) {
	if page < 1 || limit < 1 || limit > 25 {
		return nil, nil, -1, commonerrors.NewInvalidInput("invalid page or limit")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imagesMetadata, totalCount, err := s.imagesDBRepo.GetImagesMetadataByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, nil, -1, commonerrors.NewInternal(fmt.Sprintf("error reading images metadata from database: %v", err))
	}

	imagesBytes := make([][]byte, len(imagesMetadata))
	for i, imageMetadata := range imagesMetadata {
		previewImageObjectName := domain.CreatePreviewImageObjectName(imageMetadata.ID)
		imageBytes, err := s.imagesCacheRepo.GetImage(ctx, previewImageObjectName)
		if err != nil {
			return nil, nil, -1, commonerrors.NewInternal(fmt.Sprintf("error reading image from cache: %v", err))
		}
		if imageBytes == nil {
			imageBytes, err = s.imagesStorageRepo.DownloadImage(ctx, previewImageObjectName)
			if err != nil {
				return nil, nil, -1, commonerrors.NewInternal(fmt.Sprintf("error downloading image from storage: %v", err))
			}
			err = s.imagesCacheRepo.CacheImage(ctx, previewImageObjectName, imageBytes, s.cacheExpiry)
			if err != nil {
				return nil, nil, -1, commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
			}
		}
		imagesBytes[i] = imageBytes
	}

	return imagesMetadata, imagesBytes, totalCount, nil
}

func (s *ImagesService) UpdateDetails(userID uuid.UUID, oldName, newName, newDescription string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageMetadata, err := s.imagesDBRepo.GetImageMetadataByUserIDAndName(ctx, userID, oldName)
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

	err = s.imagesDBRepo.UpdateImageMetadataDetails(ctx, imageMetadata.ID, newName, newDescription)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error updating image metadata in database: %v", err))
	}

	return nil
}

func (s *ImagesService) Transform(
	userID uuid.UUID,
	name string,
	transformations []domain.Transformation,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageMetadata, err := s.imagesDBRepo.GetImageMetadataByUserIDAndName(ctx, userID, name)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error reading image metadata from database: %v", err))
	}

	fullImageObjectName := domain.CreateFullImageObjectName(imageMetadata.ID)
	imageBytes, err := s.imagesCacheRepo.GetImage(ctx, fullImageObjectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error reading image from cache: %v", err))
	}
	if imageBytes == nil {
		imageBytes, err = s.imagesStorageRepo.DownloadImage(ctx, fullImageObjectName)
		if err != nil {
			return commonerrors.NewInternal(fmt.Sprintf("error downloading image from storage: %v", err))
		}
		err = s.imagesCacheRepo.CacheImage(ctx, fullImageObjectName, imageBytes, s.cacheExpiry)
		if err != nil {
			return commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
		}
	}

	transformedBytes, err := s.transformationsService.Apply(imageBytes, transformations)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error applying transformations: %v", err))
	}

	err = s.imagesStorageRepo.UploadImage(ctx, fullImageObjectName, transformedBytes)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error uploading image to storage: %v", err))
	}
	err = s.imagesCacheRepo.CacheImage(ctx, fullImageObjectName, transformedBytes, s.cacheExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error caching image: %v", err))
	}

	previewBytes, err := s.transformationsService.CreatePreview(transformedBytes)
	previewImageObjectName := domain.CreatePreviewImageObjectName(imageMetadata.ID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error creating preview image: %v", err))
	}
	err = s.imagesStorageRepo.UploadImage(ctx, previewImageObjectName, previewBytes)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error uploading preview image to storage: %v", err))
	}
	err = s.imagesCacheRepo.CacheImage(ctx, previewImageObjectName, previewBytes, s.cacheExpiry)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error caching preview image: %v", err))
	}

	err = s.imagesDBRepo.UpdateImageMetadataUpdatedAt(ctx, imageMetadata.ID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error updating image metadata in database: %v", err))
	}

	return nil
}

func (s *ImagesService) Delete(userID uuid.UUID, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageMetadata, err := s.imagesDBRepo.GetImageMetadataByUserIDAndName(ctx, userID, name)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error reading image metadata from database: %v", err))
	}

	err = s.imagesDBRepo.DeleteImageMetadata(ctx, imageMetadata.ID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting image metadata from database: %v", err))
	}

	fullImageObjectName := domain.CreateFullImageObjectName(imageMetadata.ID)
	err = s.imagesStorageRepo.DeleteImage(ctx, fullImageObjectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting image from storage: %v", err))
	}

	err = s.imagesCacheRepo.DeleteImage(ctx, fullImageObjectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting image from cache: %v", err))
	}

	previewImageObjectName := domain.CreatePreviewImageObjectName(imageMetadata.ID)
	err = s.imagesStorageRepo.DeleteImage(ctx, previewImageObjectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting preview image from storage: %v", err))
	}

	err = s.imagesCacheRepo.DeleteImage(ctx, previewImageObjectName)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting preview image from cache: %v", err))
	}

	return nil
}

func (s *ImagesService) AdminListAllImages(page, limit int) ([]*domain.ImageMetadata, int, error) {
	if page < 1 || limit < 1 || limit > 25 {
		return nil, -1, commonerrors.NewInvalidInput("invalid page or limit")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	images, total, err := s.imagesDBRepo.GetAllImagesMetadata(ctx, page, limit)
	if err != nil {
		return nil, 0, commonerrors.NewInternal(fmt.Sprintf("error fetching images from database: %v", err))
	}

	return images, total, nil
}

func (s *ImagesService) AdminDeleteImage(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.imagesDBRepo.DeleteImageMetadata(ctx, id)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error deleting image from database: %v", err))
	}

	return nil
}
