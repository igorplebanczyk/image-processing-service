package worker

import (
	"context"
	"image-processing-service/src/internal/common/storage"
	"log/slog"
	"time"
)

type Worker struct {
	storage  *storage.Service
	ctx      context.Context
	interval time.Duration
	stop     chan bool
}

func New(storage *storage.Service) *Worker {
	return &Worker{
		storage:  storage,
		ctx:      context.Background(),
		interval: time.Hour * 24,
		stop:     make(chan bool),
	}
}

func (s *Worker) Start() {
	slog.Info("Starting storage worker")
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)

			err := s.repo.DeleteExpiredRefreshTokens(ctx)
			if err != nil {
				slog.Error("Storage error: error deleting expired refresh tokens", "error", err)
			} else {
				slog.Info("Dangling images deleted")
			}

			cancel()
		case <-s.stop:
			return
		}
	}
}

func (s *Worker) Stop() {
	s.stop <- true
	slog.Info("Storage worker stopped")
}
