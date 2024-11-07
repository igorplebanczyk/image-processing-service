package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Service struct {
	client *redis.Client
	ctx    context.Context
}

func NewService(addr string, password string, db int) (*Service, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Service{
		client: client,
		ctx:    context.Background(),
	}, nil
}

func (s *Service) Set(key string, value []byte, expiration time.Duration) error {
	err := s.client.Set(s.ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}

	return nil
}

func (s *Service) Get(key string) ([]byte, error) {
	val, err := s.client.Get(s.ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get value of key: %w", err)
	}

	return []byte(val), nil
}

func (s *Service) Delete(key string) error {
	err := s.client.Del(s.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}
