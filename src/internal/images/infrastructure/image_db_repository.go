package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/src/internal/common/database/tx"
	"image-processing-service/src/internal/images/domain"
	"log/slog"
	"time"
)

type ImageDBRepository struct {
	db         *sql.DB
	txProvider *tx.Provider
}

func NewImageDBRepository(db *sql.DB, txProvider *tx.Provider) *ImageDBRepository {
	return &ImageDBRepository{db: db, txProvider: txProvider}
}

func (r *ImageDBRepository) CreateImageMetadata(
	ctx context.Context,
	userID uuid.UUID,
	name,
	description string,
) (*domain.ImageMetadata, error) {
	slog.Info("DB query", "operation", "INSERT", "table", "images_metadata", "parameters", fmt.Sprintf("userID: %s, name: %s, description: %s", userID, name, description))

	imageMetadata := domain.NewImageMetadata(userID, name, description)
	err := r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO images_metadata (id, user_id, name, description, created_at, updated_at) 
											VALUES ($1, $2, $3, $4, $5, $6)`,
			imageMetadata.ID, imageMetadata.UserID, imageMetadata.Name, imageMetadata.Description, imageMetadata.CreatedAt, imageMetadata.UpdatedAt)
		if err != nil {
			return fmt.Errorf("error creating image metadata: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error creating image metadata: %w", err)
	}

	return imageMetadata, nil
}

func (r *ImageDBRepository) GetImageMetadataByUserIDAndName(ctx context.Context, userID uuid.UUID, name string) (*domain.ImageMetadata, error) {
	slog.Info("DB query", "operation", "SELECT", "table", "images_metadata", "parameters", fmt.Sprintf("userID: %s, name: %s", userID, name))

	var imageMetadata domain.ImageMetadata
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, name, description, created_at, updated_at 
										FROM images_metadata 
										WHERE user_id = $1 AND name = $2`, userID, name)
	err := row.Scan(&imageMetadata.ID, &imageMetadata.UserID, &imageMetadata.Name, &imageMetadata.Description, &imageMetadata.CreatedAt, &imageMetadata.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting image metadata: %w", err)
	}

	return &imageMetadata, nil
}

func (r *ImageDBRepository) GetImagesMetadataByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*domain.ImageMetadata, int, error) {
	slog.Info("DB query", "operation", "SELECT", "table", "images_metadata", "parameters", fmt.Sprintf("userID: %s, page: %d, limit: %d", userID, page, limit))

	offset := (page - 1) * limit

	var imagesMetadata []*domain.ImageMetadata
	var total int

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, name, description, created_at, updated_at 
										FROM images_metadata 
										WHERE user_id = $1
										ORDER BY created_at DESC
										LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, -1, fmt.Errorf("error getting images metadata: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var imageMetadata domain.ImageMetadata
		err := rows.Scan(&imageMetadata.ID, &imageMetadata.UserID, &imageMetadata.Name, &imageMetadata.Description, &imageMetadata.CreatedAt, &imageMetadata.UpdatedAt)
		if err != nil {
			return nil, -1, fmt.Errorf("error scanning image metadata: %w", err)
		}

		imagesMetadata = append(imagesMetadata, &imageMetadata)
	}

	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM images_metadata WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, -1, fmt.Errorf("error getting total images metadata: %w", err)
	}

	return imagesMetadata, total, nil
}

func (r *ImageDBRepository) GetAllImagesMetadata(ctx context.Context, page, limit int) ([]*domain.ImageMetadata, int, error) {
	slog.Info("DB query", "operation", "SELECT", "table", "images_metadata", "parameters", fmt.Sprintf("page: %d, limit: %d", page, limit))

	offset := (page - 1) * limit

	var imagesMetadata []*domain.ImageMetadata
	var total int

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, name, description, created_at, updated_at
										FROM images_metadata
										ORDER BY created_at DESC
										LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, -1, fmt.Errorf("error getting all images metadata: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var imageMetadata domain.ImageMetadata
		err := rows.Scan(&imageMetadata.ID, &imageMetadata.UserID, &imageMetadata.Name, &imageMetadata.Description, &imageMetadata.CreatedAt, &imageMetadata.UpdatedAt)
		if err != nil {
			return nil, -1, fmt.Errorf("error scanning image metadata: %w", err)
		}

		imagesMetadata = append(imagesMetadata, &imageMetadata)
	}

	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM images_metadata`).Scan(&total)
	if err != nil {
		return nil, -1, fmt.Errorf("error getting total images metadata: %w", err)
	}

	return imagesMetadata, total, nil
}

func (r *ImageDBRepository) UpdateImageMetadataDetails(ctx context.Context, id uuid.UUID, newName, newDescription string) error {
	slog.Info("DB query", "operation", "UPDATE", "table", "images_metadata", "parameters", fmt.Sprintf("id: %s, newName: %s, newDescription: %s", id, newName, newDescription))

	err := r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE images_metadata 
										SET name = $1, description = $2, updated_at = $3 
										WHERE id = $4`, newName, newDescription, time.Now(), id)
		if err != nil {
			return fmt.Errorf("error updating image metadata details: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error updating image metadata details: %w", err)
	}

	return nil
}

func (r *ImageDBRepository) UpdateImageMetadataUpdatedAt(ctx context.Context, id uuid.UUID) error {
	slog.Info("DB query", "operation", "UPDATE", "table", "images_metadata", "parameters", fmt.Sprintf("id: %s", id))

	err := r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE images_metadata SET updated_at = $1 WHERE id = $2`, time.Now(), id)
		if err != nil {
			return fmt.Errorf("error updating image metadata updated at: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error updating image metadata updated at: %w", err)
	}

	return nil
}

func (r *ImageDBRepository) DeleteImageMetadata(ctx context.Context, id uuid.UUID) error {
	slog.Info("DB query", "operation", "DELETE", "table", "images_metadata", "parameters", fmt.Sprintf("id: %s", id))

	err := r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM images_metadata WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("error deleting image metadata: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting image metadata: %w", err)
	}

	return nil
}
