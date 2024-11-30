package storage

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"log/slog"
)

type Service struct {
	client        *azblob.Client
	containerName string
}

func NewService(accountName, accountKey, serviceURL, containerName string) (*Service, error) {
	slog.Info("Init step 9: connecting to storage...", "account", accountName, "container", containerName)
	key, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("error creating shared key credential: %w", err)
	}

	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, key, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	slog.Info("Init step 10: connected to storage")

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
