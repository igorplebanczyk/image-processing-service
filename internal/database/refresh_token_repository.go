package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/auth"
	"time"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) (*auth.RefreshToken, error) {
	id := uuid.New()
	createdAt := time.Now()

	refreshToken := &auth.RefreshToken{
		ID:        id,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}

	_, err := r.db.Exec(`INSERT INTO refresh_tokens(id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`,
		refreshToken.ID, refreshToken.UserID, refreshToken.Token, refreshToken.ExpiresAt, refreshToken.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating refresh token: %w", err)
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepository) GetRefreshTokenByValue(field, value string) (*auth.RefreshToken, error) {
	var refreshToken auth.RefreshToken

	query := fmt.Sprintf(`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE %s = $1`, field)
	row := r.db.QueryRow(query, value)
	err := row.Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(userID uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("error revoking refresh token: %w", err)
	}

	return nil
}
