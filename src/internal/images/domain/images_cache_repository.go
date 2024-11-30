package domain

import (
	"context"
	"time"
)

type ImagesCacheRepository interface {
	CacheImage(ctx context.Context, key string, bytes []byte, expiry time.Duration) error
	GetImage(ctx context.Context, key string) ([]byte, error)
	DeleteImage(ctx context.Context, key string) error
}
