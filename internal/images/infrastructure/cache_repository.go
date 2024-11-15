package infrastructure

import (
	"context"
	"fmt"
	"image-processing-service/internal/common/cache"
	"log/slog"
	"time"
)

type ImageCacheRepository struct {
	cache *cache.Service
}

func NewImageCacheRepository(cache *cache.Service) *ImageCacheRepository {
	return &ImageCacheRepository{cache: cache}
}

func (r *ImageCacheRepository) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	slog.Info("Setting key in cache", "key", key)

	err := r.cache.Client().Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}

	return nil
}

func (r *ImageCacheRepository) Get(ctx context.Context, key string) ([]byte, error) {
	slog.Info("Getting key from cache", "key", key)

	val, err := r.cache.Client().Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get value of key: %w", err)
	}

	return []byte(val), nil
}

func (r *ImageCacheRepository) Delete(ctx context.Context, key string) error {
	slog.Info("Deleting key from cache", "key", key)

	err := r.cache.Client().Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}
