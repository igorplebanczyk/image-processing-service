package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/services/auth/util"
	"image-processing-service/internal/users"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, username, email, password string) (*users.User, error) {
	id := uuid.New()

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	createdAt := time.Now()

	user := &users.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	_, err = r.db.ExecContext(ctx, `INSERT INTO users (id, username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating users: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*users.User, error) {
	var user users.User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting users by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*users.User, error) {
	var user users.User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting users by username: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	var user users.User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting users by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, username, email string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE users SET username = $1, email = $2, updated_at = $3 WHERE id = $4`, username, email, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error updating users: %w", err)
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting users: %w", err)
	}

	return nil
}
