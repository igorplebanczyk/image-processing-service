package database

import (
	"database/sql"
	"errors"
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

func (r *UserRepository) CreateUser(username, email, password string) (*users.User, error) {
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

	_, err = r.db.Exec(`INSERT INTO users (id, username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating users: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByValue(field, value string) (*users.User, error) {
	var user users.User

	query := fmt.Sprintf(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE %s = $1`, field)
	row := r.db.QueryRow(query, value)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("users with %s %s not found", field, value)
		}
		return nil, fmt.Errorf("error getting users by %s: %w", field, err)
	}

	return &user, nil
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting users: %w", err)
	}

	return nil
}

func (r *UserRepository) UpdateUser(id uuid.UUID, username, email string) error {
	_, err := r.db.Exec(`UPDATE users SET username = $1, email = $2, updated_at = $3 WHERE id = $4`, username, email, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error updating users: %w", err)
	}

	return nil
}
