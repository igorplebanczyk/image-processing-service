package worker

import (
	"context"
	"database/sql"
	"fmt"
	"image-processing-service/src/internal/common/storage"
	"log/slog"
	"time"
)

// Worker is a storage worker that periodically deletes dangling images from the storage.
// Dangling images are images that are stored in the storage but are not present in the database.
// This can happen for example when a user deletes their account, in which case the images are deleted from the database
// automatically (via ON DELETE CASCADE), but not from the storage. I could have done that, but for simplicity I decided
// to implement this worker instead. Another benefit of this is that if a database is wiped (which I did multiple times
// during development), there is no need to manually delete the images from the storage.

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
	slog.Info("Init step 19: storage worker started")

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
	slog.Info("Shutdown step 3: worker stopped")
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
