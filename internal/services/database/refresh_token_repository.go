package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/users"
	"time"
)

type RefreshTokenRepository struct {
	service *Service
}

func NewRefreshTokenRepository(service *Service) *RefreshTokenRepository {
	return &RefreshTokenRepository{service: service}
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

	err := r.service.withTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO refresh_tokens(id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`,
			refreshToken.ID, refreshToken.UserID, refreshToken.Token, refreshToken.ExpiresAt, refreshToken.CreatedAt)
		if err != nil {
			return fmt.Errorf("error creating refresh token: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error creating refresh token: %w", err)
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepository) GetRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) ([]*users.RefreshToken, error) {
	var refreshTokens []*users.RefreshToken

	rows, err := r.service.DB.QueryContext(
		ctx,
		`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting refresh tokens by user ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var refreshToken users.RefreshToken
		err := rows.Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning refresh token: %w", err)
		}

		refreshTokens = append(refreshTokens, &refreshToken)
	}

	return refreshTokens, nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(ctx context.Context, userID uuid.UUID) error {
	return r.service.withTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
		if err != nil {
			return fmt.Errorf("error revoking refresh token: %w", err)
		}

		return nil
	})
}

func (r *RefreshTokenRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	return r.service.withTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE expires_at < $1`, time.Now())
		if err != nil {
			return fmt.Errorf("error deleting expired refresh tokens: %w", err)
		}

		return nil
	})
}
