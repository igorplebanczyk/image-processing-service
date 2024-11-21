package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"image-processing-service/src/internal/common/logs"
	"log/slog"
)

type Service struct {
	client *redis.Client
}

func NewService(host, port, password string, db int) (*Service, error) {
	slog.Info("Connecting to cache...", "type", logs.Standard)
	addr := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	slog.Info("Connected to cache", "type", logs.Standard)
	return &Service{
		client: client,
	}, nil
}

func (s *Service) Client() *redis.Client {
	return s.client
}
