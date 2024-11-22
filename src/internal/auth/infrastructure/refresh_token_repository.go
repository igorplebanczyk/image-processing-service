package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/src/internal/auth/domain"
	"image-processing-service/src/internal/common/database/transactions"
	"log/slog"
	"time"
)

type RefreshTokenRepository struct {
	db         *sql.DB
	txProvider *transactions.TransactionProvider
}

func NewRefreshTokenRepository(db *sql.DB, txProvider *transactions.TransactionProvider) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db, txProvider: txProvider}
}

func (r *RefreshTokenRepository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	slog.Info("DB query", "operation", "INSERT", "table", "refresh_tokens", "parameters", fmt.Sprintf("user_id: %s, token: %s, expires_at: %s", userID, token, expiresAt))

	id := uuid.New()
	createdAt := time.Now()

	return r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO refresh_tokens(id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`,
			id, userID, token, expiresAt, createdAt)
		if err != nil {
			return fmt.Errorf("error creating refresh token: %w", err)
		}

		return nil
	})
}

func (r *RefreshTokenRepository) GetRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.RefreshToken, error) {
	slog.Info("DB query", "operation", "SELECT", "table", "refresh_tokens", "parameters", fmt.Sprintf("user_id: %s", userID))

	var refreshTokens []*domain.RefreshToken

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting refresh tokens by user ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var refreshToken domain.RefreshToken
		err := rows.Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning refresh token: %w", err)
		}

		refreshTokens = append(refreshTokens, &refreshToken)
	}

	return refreshTokens, nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(ctx context.Context, userID uuid.UUID) error {
	slog.Info("DB query", "operation", "DELETE", "table", "refresh_tokens", "parameters", fmt.Sprintf("user_id: %s", userID))

	return r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
		if err != nil {
			return fmt.Errorf("error revoking refresh token: %w", err)
		}

		return nil
	})
}

func (r *RefreshTokenRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	slog.Info("DB query", "operation", "DELETE", "table", "refresh_tokens", "parameters", fmt.Sprintf("user_id: %s", userID))

	return r.txProvider.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
		if err != nil {
			return fmt.Errorf("error revoking refresh tokens: %w", err)
		}

		return nil
	})
}
