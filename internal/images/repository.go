package images

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Repository struct {
	db      *sql.DB
	storage StorageService
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateImage(ctx context.Context, userID uuid.UUID, name string, image []byte) error {
	url, err := r.storage.UploadObject(ctx, userID.String(), name, image)
	if err != nil {
		return fmt.Errorf("error uploading image: %w", err)
	}

	id := uuid.New()
	createdAt := time.Now()

	_, err = r.db.Exec(`INSERT INTO images(id, user_id, name, url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5 $6)`,
		id, userID, name, url, createdAt, createdAt)
	if err != nil {
		return fmt.Errorf("error creating image: %w", err)
	}

	return nil
}
