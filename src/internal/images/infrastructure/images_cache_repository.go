package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"image-processing-service/src/internal/common/cache"
	"image-processing-service/src/internal/common/metrics"
	"log/slog"
	"time"
)

type ImagesCacheRepository struct {
	cache *cache.Service
}

func NewImagesCacheRepository(cache *cache.Service) *ImagesCacheRepository {
	return &ImagesCacheRepository{cache: cache}
}

func (r *ImagesCacheRepository) CacheImage(ctx context.Context, key string, bytes []byte, expiry time.Duration) error {
	slog.Info("Setting key in cache", "key", key)
	metrics.CacheOperationsTotal.WithLabelValues("set").Inc()

	err := r.cache.Client().Set(ctx, key, bytes, expiry).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}

	return nil
}

func (r *ImagesCacheRepository) GetImage(ctx context.Context, key string) ([]byte, error) {
	slog.Info("Getting key from cache", "key", key)
	metrics.CacheOperationsTotal.WithLabelValues("get").Inc()

	val, err := r.cache.Client().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, cache.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	return []byte(val), nil
}

func (r *ImagesCacheRepository) DeleteImage(ctx context.Context, key string) error {
	slog.Info("Deleting key from cache", "key", key)
	metrics.CacheOperationsTotal.WithLabelValues("delete").Inc()

	err := r.cache.Client().Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}
