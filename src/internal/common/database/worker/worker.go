package worker

import (
	"context"
	"database/sql"
	"image-processing-service/src/internal/common/database/transactions"
	"log/slog"
	"time"
)

const interval = time.Hour

type Worker struct {
	repo     repository
	ctx      context.Context
	interval time.Duration
	stop     chan bool
}

func New(db *sql.DB, txProvider *transactions.TransactionProvider) *Worker {
	return &Worker{
		repo:     newRepository(db, txProvider),
		ctx:      context.Background(),
		interval: interval,
		stop:     make(chan bool),
	}
}

func (s *Worker) Start() {
	slog.Info("Starting database worker")
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)

			err := s.repo.DeleteExpiredRefreshTokens(ctx)
			if err != nil {
				slog.Error("Database error: error deleting expired refresh tokens", "error", err)
			} else {
				slog.Info("Expired refresh tokens deleted")
			}

			cancel()
		case <-s.stop:
			return
		}
	}
}

func (s *Worker) Stop() {
	s.stop <- true
	slog.Info("Database worker stopped")
}
