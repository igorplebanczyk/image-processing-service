package images

import (
	"context"
	"time"
)

type StorageService interface {
	UploadObject(ctx context.Context, blobName string, data []byte) error
	DownloadObject(ctx context.Context, blobName string) ([]byte, error)
	DeleteObject(ctx context.Context, blobName string) error
}

type CacheService interface {
	Set(key string, value []byte, expiration time.Duration) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}
