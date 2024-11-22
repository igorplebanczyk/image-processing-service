package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log/slog"
)

type Service struct {
	client *redis.Client
}

func NewService(host, port, password string, db int) (*Service, error) {
	slog.Info("Connecting to cache...", "host", host, "port", port, "db", db)
	addr := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	slog.Info("Connected to cache")
	return &Service{
		client: client,
	}, nil
}

func (s *Service) Client() *redis.Client {
	return s.client
}
