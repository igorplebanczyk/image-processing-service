package infrastructure

import (
	"context"
	"fmt"
	"image-processing-service/src/internal/common/metrics"
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

func (r *ImageStorageRepository) UploadImage(ctx context.Context, name string, bytes []byte) error {
	slog.Info("Uploading file to storage", "blob_name", name)
	metrics.StorageOperationsTotal.WithLabelValues("upload").Inc()

	_, err := r.storage.Client().UploadBuffer(ctx, r.storage.ContainerName(), name, bytes, nil)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	return nil
}

func (r *ImageStorageRepository) DownloadImage(ctx context.Context, name string) ([]byte, error) {
	slog.Info("Downloading file from storage", "blob_name", name)
	metrics.StorageOperationsTotal.WithLabelValues("download").Inc()

	var data = make([]byte, domain.MaxImageSize)
	_, err := r.storage.Client().DownloadBuffer(ctx, r.storage.ContainerName(), name, data, nil)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %w", err)
	}

	return data, nil
}

func (r *ImageStorageRepository) DeleteImage(ctx context.Context, name string) error {
	slog.Info("Deleting file from storage", "blob_name", name)
	metrics.StorageOperationsTotal.WithLabelValues("delete").Inc()

	_, err := r.storage.Client().DeleteBlob(ctx, r.storage.ContainerName(), name, nil)
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}

	return nil
}
