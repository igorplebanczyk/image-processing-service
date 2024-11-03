package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/auth"
	"time"
)

const passwordRegex = `(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,32}`

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(username, email, password string) (*auth.User, error) {
	id := uuid.New()

	err := validate(r, username, email, password)
	if err != nil {
		return nil, fmt.Errorf("error validating user: %w", err)
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	createdAt := time.Now()

	user := &auth.User{
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
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func GetUserByID(id uuid.UUID) (*auth.User, error) {

	return nil, nil
}

func GetUserByUsername(username string) (*auth.User, error) {

	return nil, nil
}

func GetUserByEmail(email string) (*auth.User, error) {

	return nil, nil
}
