package images

import (
	"context"
	"time"
)

type StorageService interface {
	Upload(ctx context.Context, objectName string, data []byte) error
	Download(ctx context.Context, objectName string) ([]byte, error)
	Delete(ctx context.Context, objectName string) error
}

type CacheService interface {
	Set(key string, value []byte, expiration time.Duration) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type TransformationService interface {
	Transform(data []byte, transformation string, options struct{ any }) ([]byte, error)
}
