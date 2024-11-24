package worker

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type imagesRepository struct {
	db *sql.DB
}

func newImagesRepository(db *sql.DB) imagesRepository {
	return imagesRepository{
		db: db,
	}
}

type image struct {
	name   string
	userID uuid.UUID
}

func (r *imagesRepository) getAllImages(ctx context.Context) ([]image, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT name, user_id FROM images")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []image
	for rows.Next() {
		var img image
		err = rows.Scan(&img.name, &img.userID)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return images, nil
}
