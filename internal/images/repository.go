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
	url, err := r.storage.UploadObject(ctx, userID.String(), image)
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

func (r *Repository) GetImagesByUserID(userID uuid.UUID) ([]Image, error) {
	var images []Image

	rows, err := r.db.Query(`SELECT id, name, url, created_at, updated_at FROM images WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting images by user id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var image Image
		err := rows.Scan(&image.ID, &image.Name, &image.URL, &image.CreatedAt, &image.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning image: %w", err)
		}

		images = append(images, image)
	}

	return images, nil
}

func (r *Repository) GetImageByID(ctx context.Context, id uuid.UUID) ([]byte, error) {
	var image Image

	row := r.db.QueryRow(`SELECT id, name, url, created_at, updated_at FROM images WHERE id = $1`, id)
	err := row.Scan(&image.ID, &image.Name, &image.URL, &image.CreatedAt, &image.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting image by id: %w", err)
	}

	imageBytes, err := r.storage.DownloadObject(ctx, image.URL)
	if err != nil {
		return nil, fmt.Errorf("error downloading image: %w", err)
	}

	return imageBytes, nil
}
