package worker

import (
	"log/slog"
	"time"
)

type Service struct {
	repo     RefreshTokenRepository
	interval time.Duration
	stop     chan bool
}

func New(repo RefreshTokenRepository, interval time.Duration) *Service {
	return &Service{
		repo:     repo,
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
			err := s.repo.DeleteExpiredRefreshTokens()
			if err != nil {
				slog.Error("error deleting expired refresh tokens", "error", err)
			} else {
				slog.Info("expired refresh tokens deleted")
			}
		case <-s.stop:
			slog.Info("Worker stopped")
			return
		}
	}
}

func (s *Service) Stop() {
	s.stop <- true
}
