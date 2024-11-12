package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/common/database/transactions"
	"image-processing-service/internal/users/domain"
)

type UserRepository struct {
	db         *sql.DB
	txProvider *transactions.TransactionProvider
}

func NewUserRepository(db *sql.DB, txProvider *transactions.TransactionProvider) *UserRepository {
	return &UserRepository{db: db, txProvider: txProvider}
}

func (r *UserRepository) CreateUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	user := domain.NewUser(username, email, password)

	err := r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO users (id, username, email, password, created_at, updated_at) 
											VALUES ($1, $2, $3, $4, $5, $6)`,
			user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
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
	var user domain.User

	err := r.db.QueryRowContext(ctx, `SELECT * FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, username, email string) error {
	return r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE users SET username = $1, email = $2 WHERE id = $3`, username, email, id)
		if err != nil {
			return fmt.Errorf("error updating user: %w", err)
		}

		return nil
	})
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("error deleting user: %w", err)
		}

		return nil
	})
}
