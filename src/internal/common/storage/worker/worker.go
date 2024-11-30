package worker

import (
	"context"
	"database/sql"
	"fmt"
	"image-processing-service/src/internal/common/storage"
	"log/slog"
	"time"
)

const interval = 24 * time.Hour

type Worker struct {
	repo     imagesDBRepository
	storage  imagesStorageRepository
	ctx      context.Context
	stop     context.CancelFunc
	interval time.Duration
}

func New(db *sql.DB, storage *storage.Service) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		repo:     newImagesDBRepository(db),
		storage:  newImagesStorageRepository(storage),
		ctx:      ctx,
		interval: interval,
		stop:     cancel,
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

			err := s.deleteDanglingImages(ctx)
			if err != nil {
				slog.Error("error deleting dangling images", "error", err)
			} else {
				slog.Info("Dangling images deleted")
			}

			cancel()
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Worker) Stop() {
	s.stop()
	slog.Info("Storage worker stopped")
}

func (s *Worker) deleteDanglingImages(ctx context.Context) error {
	imagesNamesStorage, err := s.storage.getAllImagesNames()
	if err != nil {
		return fmt.Errorf("failed to get images names from storage: %w", err)
	}

	imagesDB, err := s.repo.getAllImages(ctx)
	if err != nil {
		return fmt.Errorf("failed to get images from database: %w", err)
	}

	danglingImages, err := getAllDanglingImagesNames(imagesNamesStorage, imagesDB)
	if err != nil {
		return fmt.Errorf("failed to delete dangling images: %w", err)
	}

	for _, name := range danglingImages {
		err := s.storage.deleteImage(name)
		if err != nil {
			return fmt.Errorf("failed to delete dangling image %s: %w", name, err)
		}
	}

	return nil
}
