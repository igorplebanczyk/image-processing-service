package infrastructure

import (
	"context"
	"fmt"
	"image-processing-service/src/internal/common/storage"
	"image-processing-service/src/internal/images/domain"
	"log/slog"
)

type ImageStorageRepository struct {
	storage *storage.Service
}

func NewImageStorageRepository(storage *storage.Service) *ImageStorageRepository {
	return &ImageStorageRepository{storage: storage}
}

func (r *ImageStorageRepository) Upload(ctx context.Context, blobName string, data []byte) error {
	slog.Info("Uploading file to storage", "blob_name", blobName)

	_, err := r.storage.Client().UploadBuffer(ctx, r.storage.ContainerName(), blobName, data, nil)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	return nil
}

func (r *ImageStorageRepository) Download(ctx context.Context, blobName string) ([]byte, error) {
	slog.Info("Downloading file from storage", "blob_name", blobName)

	var data = make([]byte, domain.MaxImageSize)
	_, err := r.storage.Client().DownloadBuffer(ctx, r.storage.ContainerName(), blobName, data, nil)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %w", err)
	}

	return data, nil
}

func (r *ImageStorageRepository) Delete(ctx context.Context, blobName string) error {
	slog.Info("Deleting file from storage", "blob_name", blobName)

	_, err := r.storage.Client().DeleteBlob(ctx, r.storage.ContainerName(), blobName, nil)
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}

	return nil
}
