package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
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
	slog.Info("DB query", "operation", "SELECT", "table", "users", "parameters", fmt.Sprintf("username: %s", username))

	var user domain.User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, password, role FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("error getting user by username: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserRoleByID(ctx context.Context, id uuid.UUID) (domain.Role, error) {
	slog.Info("DB query", "operation", "SELECT", "table", "users", "parameters", fmt.Sprintf("id: %s", id))

	var role domain.Role

	err := r.db.QueryRowContext(
		ctx,
		`SELECT role FROM users WHERE id = $1`,
		id,
	).Scan(&role)
	if err != nil {
		return "", fmt.Errorf("error getting user role by ID: %w", err)
	}

	return role, nil
}
