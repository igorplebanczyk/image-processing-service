package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/images"
	"time"
)

type ImageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) CreateImage(userID uuid.UUID, name string) (*images.Image, error) {
	id := uuid.New()

	image := &images.Image{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.db.Exec(`INSERT INTO images (id, user_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		image.ID, image.UserID, image.Name, image.CreatedAt, image.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating image: %w", err)
	}

	return image, nil
}

func (r *ImageRepository) GetImagesByUserID(userID uuid.UUID) ([]*images.Image, error) {
	var imagesList []*images.Image

	rows, err := r.db.Query(`SELECT id, user_id, name, created_at, updated_at FROM images WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting images by user id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var image images.Image
		err := rows.Scan(&image.ID, &image.UserID, &image.Name, &image.CreatedAt, &image.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning image: %w", err)
		}
		imagesList = append(imagesList, &image)
	}

	return imagesList, nil
}

func (r *ImageRepository) GetImageByUserIDandName(userID uuid.UUID, name string) (*images.Image, error) {
	var img images.Image

	row := r.db.QueryRow(`SELECT id, user_id, name, created_at, updated_at FROM images WHERE user_id = $1 AND name = $2`, userID, name)

	err := row.Scan(&img.ID, &img.UserID, &img.Name, &img.CreatedAt, &img.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error scanning image: %w", err)
	}

	return &img, nil
}

func (r *ImageRepository) UpdateImage(id uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE images SET updated_at = $1 WHERE id = $2`, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error updating image: %w", err)
	}

	return nil
}

func (r *ImageRepository) DeleteImage(id uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM images WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting image: %w", err)
	}

	return nil
}
