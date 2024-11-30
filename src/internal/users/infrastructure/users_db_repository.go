package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/src/internal/common/database/tx"
	"image-processing-service/src/internal/common/metrics"
	"image-processing-service/src/internal/users/domain"
	"log/slog"
	"time"
)

type UsersDBRepository struct {
	db         *sql.DB
	txProvider *tx.Provider
}

func NewUsersDBRepository(db *sql.DB, txProvider *tx.Provider) *UsersDBRepository {
	return &UsersDBRepository{db: db, txProvider: txProvider}
}

func (r *UsersDBRepository) CreateUser(ctx context.Context, username, email, password, otpSecret string) (*domain.User, error) {
	slog.Info("DB query", "operation", "INSERT", "table", "users", "parameters", fmt.Sprintf("username: %s, email: %s, password: %s", username, email, password))
	metrics.DBQueriesTotal.WithLabelValues("INSERT").Inc()

	user := domain.NewUser(username, email, password, otpSecret)

	err := r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO users (id, username, email, password, role, verified, otp_secret, created_at, updated_at) 
											VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			user.ID, user.Username, user.Email, user.Password, user.Role, user.Verified, user.OTPSecret, user.CreatedAt, user.UpdatedAt)
		if err != nil {
			return fmt.Errorf("error creating users: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func (r *UsersDBRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	slog.Info("DB query", "operation", "SELECT", "table", "users", "parameters", fmt.Sprintf("id: %s", id))
	metrics.DBQueriesTotal.WithLabelValues("SELECT").Inc()

	var user domain.User

	err := r.db.QueryRowContext(ctx, `SELECT * FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.OTPSecret, &user.Verified, &user.UpdatedAt, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

func (r *UsersDBRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	slog.Info("DB query", "operation", "SELECT", "table", "users", "parameters", fmt.Sprintf("email: %s", email))
	metrics.DBQueriesTotal.WithLabelValues("SELECT").Inc()

	var user domain.User

	err := r.db.QueryRowContext(ctx, `SELECT * FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.OTPSecret, &user.Verified, &user.UpdatedAt, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

func (r *UsersDBRepository) GetAllUsers(ctx context.Context, page, limit int) ([]domain.User, int, error) {
	slog.Info("DB query", "operation", "SELECT", "table", "users")
	metrics.DBQueriesTotal.WithLabelValues("SELECT").Inc()

	offset := (page - 1) * limit

	var total int
	var users []domain.User

	rows, err := r.db.QueryContext(ctx, `SELECT * FROM users LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, -1, fmt.Errorf("error getting users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.OTPSecret, &user.Verified, &user.UpdatedAt, &user.CreatedAt)
		if err != nil {
			return nil, -1, fmt.Errorf("error scanning user: %w", err)
		}

		users = append(users, user)
	}

	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, -1, fmt.Errorf("error getting total users: %w", err)
	}

	return users, total, nil
}

func (r *UsersDBRepository) UpdateUserDetails(ctx context.Context, id uuid.UUID, username, email string) error {
	slog.Info("DB query", "operation", "UPDATE", "table", "users", "parameters", fmt.Sprintf("id: %s, username: %s, email: %s", id, username, email))
	metrics.DBQueriesTotal.WithLabelValues("UPDATE").Inc()

	return r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE users SET username = $1, email = $2, updated_at = $3 WHERE id = $4`,
			username, email, time.Now(), id)
		if err != nil {
			return fmt.Errorf("error updating user: %w", err)
		}

		return nil
	})
}

func (r *UsersDBRepository) UpdateUserRole(ctx context.Context, id uuid.UUID, role domain.Role) error {
	slog.Info("DB query", "operation", "UPDATE", "table", "users", "parameters", fmt.Sprintf("id: %s, role: %s", id, role))
	metrics.DBQueriesTotal.WithLabelValues("UPDATE").Inc()

	return r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE users SET role = $1, updated_at = $2 WHERE id = $3`, role, time.Now(), id)
		if err != nil {
			return fmt.Errorf("error updating user role: %w", err)
		}

		return nil
	})
}

func (r *UsersDBRepository) UpdateUserAsVerified(ctx context.Context, id uuid.UUID) error {
	slog.Info("DB query", "operation", "UPDATE", "table", "users", "parameters", fmt.Sprintf("id: %s", id))
	metrics.DBQueriesTotal.WithLabelValues("UPDATE").Inc()

	return r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE users SET verified = true, updated_at = $1 WHERE id = $2`, time.Now(), id)
		if err != nil {
			return fmt.Errorf("error updating user as verified: %w", err)
		}

		return nil
	})
}

func (r *UsersDBRepository) UpdateUserPassword(ctx context.Context, id uuid.UUID, password string) error {
	slog.Info("DB query", "operation", "UPDATE", "table", "users", "parameters", fmt.Sprintf("id: %s", id))
	metrics.DBQueriesTotal.WithLabelValues("UPDATE").Inc()

	return r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE users SET password = $1, updated_at = $2 WHERE id = $3`, password, time.Now(), id)
		if err != nil {
			return fmt.Errorf("error updating user password: %w", err)
		}

		return nil
	})
}

func (r *UsersDBRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	slog.Info("DB query", "operation", "DELETE", "table", "users", "parameters", fmt.Sprintf("id: %s", id))
	metrics.DBQueriesTotal.WithLabelValues("DELETE").Inc()

	return r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("error deleting user: %w", err)
		}

		return nil
	})
}
