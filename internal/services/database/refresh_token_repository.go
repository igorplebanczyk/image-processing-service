package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/users"
	"time"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) (*users.RefreshToken, error) {
	id := uuid.New()
	createdAt := time.Now()

	refreshToken := &users.RefreshToken{
		ID:        id,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}

	_, err := r.db.ExecContext(ctx, `INSERT INTO refresh_tokens(id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`,
		refreshToken.ID, refreshToken.UserID, refreshToken.Token, refreshToken.ExpiresAt, refreshToken.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating refresh token: %w", err)
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepository) GetRefreshTokenByUserID(ctx context.Context, userID uuid.UUID) (*users.RefreshToken, error) {
	var refreshToken users.RefreshToken

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE user_id = $1`,
		userID,
	).Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("error getting refresh token by user id: %w", err)
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("error revoking refresh token: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE expires_at < $1`, time.Now())
	if err != nil {
		return fmt.Errorf("error deleting expired refresh tokens: %w", err)
	}

	return nil
}
