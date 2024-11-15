package storage

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"log/slog"
)

type Service struct {
	client        *azblob.Client
	containerName string
}

func NewService(accountName, accountKey, serviceURL, containerName string) (*Service, error) {
	key, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("error creating shared key credential: %w", err)
	}

	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, key, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	slog.Info("Connected to storage", "account", accountName, "container", containerName)

	return &Service{
		client:        client,
		containerName: containerName,
	}, nil
}

func (s *Service) Upload(ctx context.Context, blobName string, data []byte) error {
	slog.Info("Uploading file to storage", "blob_name", blobName)

	_, err := s.client.UploadBuffer(ctx, s.containerName, blobName, data, nil)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	return nil
}

func (s *Service) Download(ctx context.Context, blobName string) ([]byte, error) {
	slog.Info("Downloading file from storage", "blob_name", blobName)

	var data = make([]byte, 10*1024*1024) // 10 MB
	_, err := s.client.DownloadBuffer(ctx, s.containerName, blobName, data, nil)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %w", err)
	}

	return data, nil
}

func (s *Service) Delete(ctx context.Context, blobName string) error {
	slog.Info("Deleting file from storage", "blob_name", blobName)

	_, err := s.client.DeleteBlob(ctx, s.containerName, blobName, nil)
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}

	return nil
}
