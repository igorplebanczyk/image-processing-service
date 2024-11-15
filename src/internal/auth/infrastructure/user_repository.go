package infrastructure

import (
	"context"
	"database/sql"
	"image-processing-service/src/internal/auth/domain"
	"log/slog"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	slog.Info("DB query", "username", username)

	var user domain.User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, password FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
