package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/src/internal/common/database/transactions"
	"image-processing-service/src/internal/images/domain"
	"log/slog"
	"time"
)

type ImageRepository struct {
	db         *sql.DB
	txProvider *transactions.TransactionProvider
}

func NewImageRepository(db *sql.DB, txProvider *transactions.TransactionProvider) *ImageRepository {
	return &ImageRepository{db: db, txProvider: txProvider}
}

func (r *ImageRepository) CreateImage(ctx context.Context, userID uuid.UUID, name string) (*domain.Image, error) {
	slog.Info("DB query", "user_id", userID, "name", name)

	image := domain.NewImage(userID, name)

	err := r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO images (id, user_id, name, created_at, updated_at) 
											VALUES ($1, $2, $3, $4, $5)`,
			image.ID, image.UserID, image.Name, image.CreatedAt, image.UpdatedAt)
		if err != nil {
			return fmt.Errorf("error creating image: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error creating image: %w", err)
	}

	return image, nil
}

func (r *ImageRepository) GetImageByUserIDandName(
	ctx context.Context,
	userID uuid.UUID,
	name string,
) (*domain.Image, error) {
	slog.Info("DB query", "user_id", userID, "name", name)

	var image domain.Image

	err := r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `SELECT id, user_id, name, created_at, updated_at 
										FROM images 
										WHERE user_id = $1 AND name = $2`, userID, name)
		err := row.Scan(&image.ID, &image.UserID, &image.Name, &image.CreatedAt, &image.UpdatedAt)
		if err != nil {
			return fmt.Errorf("error getting image: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error getting image: %w", err)
	}

	return &image, nil
}

func (r *ImageRepository) GetImagesByUserID(
	ctx context.Context,
	userID uuid.UUID,
	page,
	limit *int,
) ([]*domain.Image, int, error) {
	slog.Info("DB query", "user_id", userID)

	var rows *sql.Rows
	var err error

	// If page and limit are provided, apply pagination; otherwise, fetch all results
	if page != nil && limit != nil {
		offset := (*page - 1) * (*limit)
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, user_id, name, created_at, updated_at 
			FROM images 
			WHERE user_id = $1 
			ORDER BY created_at DESC 
			LIMIT $2 OFFSET $3`, userID, *limit, offset)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, user_id, name, created_at, updated_at 
			FROM images 
			WHERE user_id = $1 
			ORDER BY created_at DESC`, userID)
	}
	if err != nil {
		return nil, -1, fmt.Errorf("error getting images by user id: %w", err)
	}
	defer rows.Close()

	var imagesList []*domain.Image
	for rows.Next() {
		var image domain.Image
		if err := rows.Scan(&image.ID, &image.UserID, &image.Name, &image.CreatedAt, &image.UpdatedAt); err != nil {
			return nil, -1, fmt.Errorf("error scanning image: %w", err)
		}
		imagesList = append(imagesList, &image)
	}

	var totalCount int
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM images WHERE user_id = $1`, userID).Scan(&totalCount)
	if err != nil {
		return nil, -1, fmt.Errorf("error getting total image count: %w", err)
	}

	return imagesList, totalCount, nil
}

func (r *ImageRepository) UpdateImage(ctx context.Context, id uuid.UUID) error {
	slog.Info("DB query", "id", id)

	return r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE images SET updated_at = $1 WHERE id = $2`, time.Now(), id)
		if err != nil {
			return fmt.Errorf("error updating image: %w", err)
		}

		return nil
	})
}

func (r *ImageRepository) DeleteImage(ctx context.Context, id uuid.UUID) error {
	slog.Info("DB query", "id", id)

	return r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM images WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("error deleting image: %w", err)
		}

		return nil
	})
}
