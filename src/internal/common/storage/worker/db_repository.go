package worker

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type imagesDBRepository struct {
	db *sql.DB
}

func newImagesDBRepository(db *sql.DB) imagesDBRepository {
	return imagesDBRepository{
		db: db,
	}
}

func (r *imagesDBRepository) getAllImages(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id FROM images_metadata")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
