package infrastructure

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
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
		`SELECT id, username, password, role FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserRoleByID(ctx context.Context, id uuid.UUID) (domain.Role, error) {
	slog.Info("DB query", "id", id)

	var role domain.Role

	err := r.db.QueryRowContext(
		ctx,
		`SELECT role FROM users WHERE id = $1`,
		id,
	).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}
