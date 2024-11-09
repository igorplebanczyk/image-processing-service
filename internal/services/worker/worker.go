package worker

import (
	"context"
	"log/slog"
	"time"
)

type Service struct {
	repo     RefreshTokenRepository
	ctx      context.Context
	interval time.Duration
	stop     chan bool
}

func New(repo RefreshTokenRepository, interval time.Duration) *Service {
	return &Service{
		repo:     repo,
		ctx:      context.Background(),
		interval: interval,
		stop:     make(chan bool),
	}
}

func (s *Service) Start() {
	slog.Info("Worker starting")
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)

			err := s.repo.DeleteExpiredRefreshTokens(ctx)
			if err != nil {
				slog.Error("error deleting expired refresh tokens", "error", err)
			} else {
				slog.Info("expired refresh tokens deleted")
			}

			cancel()
		case <-s.stop:
			slog.Info("Worker stopped")
			return
		}
	}
}

func (s *Service) Stop() {
	s.stop <- true
}
