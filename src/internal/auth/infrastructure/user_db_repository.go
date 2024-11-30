package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/src/internal/auth/domain"
	"image-processing-service/src/internal/common/metrics"
	"log/slog"
)

type UserDBRepository struct {
	db *sql.DB
}

func NewUserDBRepository(db *sql.DB) *UserDBRepository {
	return &UserDBRepository{db: db}
}

func (r *UserDBRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	slog.Info("DB query", "operation", "SELECT", "table", "users", "parameters", fmt.Sprintf("username: %s", username))
	metrics.DBQueriesTotal.WithLabelValues("SELECT").Inc()

	var user domain.User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, email, password, role, otp_secret FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.OTPSecret)
	if err != nil {
		return nil, fmt.Errorf("error getting user by username: %w", err)
	}

	return &user, nil
}

func (r *UserDBRepository) GetUserRoleByID(ctx context.Context, id uuid.UUID) (domain.Role, error) {
	slog.Info("DB query", "operation", "SELECT", "table", "users", "parameters", fmt.Sprintf("id: %s", id))
	metrics.DBQueriesTotal.WithLabelValues("SELECT").Inc()

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
