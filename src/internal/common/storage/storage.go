package storage

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"image-processing-service/src/internal/common/logs"
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

	slog.Info("Connected to storage", "type", logs.Standard)

	return &Service{
		client:        client,
		containerName: containerName,
	}, nil
}

func (s *Service) Client() *azblob.Client {
	return s.client
}

func (s *Service) ContainerName() string {
	return s.containerName
}
