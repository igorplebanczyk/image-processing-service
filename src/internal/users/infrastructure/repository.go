package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/src/internal/common/database/transactions"
	"image-processing-service/src/internal/common/logs"
	"image-processing-service/src/internal/users/domain"
	"log/slog"
	"time"
)

type UserRepository struct {
	db         *sql.DB
	txProvider *transactions.TransactionProvider
}

func NewUserRepository(db *sql.DB, txProvider *transactions.TransactionProvider) *UserRepository {
	return &UserRepository{db: db, txProvider: txProvider}
}

func (r *UserRepository) CreateUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	slog.Info("DB query", "type", logs.DB, "operation", "INSERT", "table", "users", "parameters", fmt.Sprintf("username: %s, email: %s", username, email))

	user := domain.NewUser(username, email, password)

	err := r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO users (id, username, email, password, role,created_at, updated_at) 
											VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			user.ID, user.Username, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
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

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	slog.Info("DB query", "type", logs.DB, "operation", "SELECT", "table", "users", "parameters", fmt.Sprintf("id: %s", id))

	var user domain.User

	err := r.db.QueryRowContext(ctx, `SELECT * FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserDetails(ctx context.Context, id uuid.UUID, username, email string) error {
	slog.Info("DB query", "type", logs.DB, "operation", "UPDATE", "table", "users", "parameters", fmt.Sprintf("id: %s, username: %s, email: %s", id, username, email))

	return r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE users SET username = $1, email = $2, updated_at = $3 WHERE id = $4`,
			username, email, time.Now(), id)
		if err != nil {
			return fmt.Errorf("error updating user: %w", err)
		}

		return nil
	})
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	slog.Info("DB query", "type", logs.DB, "operation", "DELETE", "table", "users", "parameters", fmt.Sprintf("id: %s", id))

	return r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("error deleting user: %w", err)
		}

		return nil
	})
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	slog.Info("DB query", "type", logs.DB, "operation", "SELECT", "table", "users")

	rows, err := r.db.QueryContext(ctx, `SELECT * FROM users`)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) UpdateUserRole(ctx context.Context, id uuid.UUID, role domain.Role) error {
	slog.Info("DB query", "type", logs.DB, "operation", "UPDATE", "table", "users", "parameters", fmt.Sprintf("id: %s, role: %s", id, role))

	return r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE users SET role = $1, updated_at = $2 WHERE id = $3`, role, time.Now(), id)
		if err != nil {
			return fmt.Errorf("error updating user role: %w", err)
		}

		return nil
	})
}
