package worker

import (
	"context"
	"database/sql"
	"image-processing-service/src/internal/common/database/tx"
	"log/slog"
	"time"
)

// Worker is a database worker that periodically deletes expired refresh tokens from the database.
// It is necessary because refresh tokens are otherwise only deleted upon logout or if an expired token is used.

const interval = time.Hour

type Worker struct {
	repo     repository
	ctx      context.Context
	stop     context.CancelFunc
	interval time.Duration
}

func New(db *sql.DB, txProvider *tx.Provider) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		repo:     newRepository(db, txProvider),
		ctx:      ctx,
		interval: interval,
		stop:     cancel,
	}
}

func (s *Worker) Start() {
	slog.Info("Init step 19: database worker started")

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
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Worker) Stop() {
	s.stop()
	slog.Info("Shutdown step 2: database worker stopped")
}
