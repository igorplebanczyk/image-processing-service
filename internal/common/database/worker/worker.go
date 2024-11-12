package worker

import (
	"context"
	"database/sql"
	"image-processing-service/internal/common/database/transactions"
	"log/slog"
	"time"
)

type Worker struct {
	repo     repository
	ctx      context.Context
	interval time.Duration
	stop     chan bool
}

func New(db *sql.DB, txProvider *transactions.TransactionProvider, interval time.Duration) *Worker {
	return &Worker{
		repo:     newRepository(db, txProvider),
		ctx:      context.Background(),
		interval: interval,
		stop:     make(chan bool),
	}
}

func (s *Worker) Start() {
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

func (s *Worker) Stop() {
	s.stop <- true
}
