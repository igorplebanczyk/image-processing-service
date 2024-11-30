package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/src/internal/auth/domain"
	"image-processing-service/src/internal/common/database/tx"
	"image-processing-service/src/internal/common/metrics"
	"log/slog"
	"time"
)

type RefreshTokenDBRepository struct {
	db         *sql.DB
	txProvider *tx.Provider
}

func NewRefreshTokenDBRepository(db *sql.DB, txProvider *tx.Provider) *RefreshTokenDBRepository {
	return &RefreshTokenDBRepository{db: db, txProvider: txProvider}
}

func (r *RefreshTokenDBRepository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	slog.Info("DB query", "operation", "INSERT", "table", "refresh_tokens", "parameters", fmt.Sprintf("user_id: %s, token: %s, expires_at: %s", userID, token, expiresAt))
	metrics.DBQueriesTotal.WithLabelValues("INSERT").Inc()

	id := uuid.New()
	createdAt := time.Now()

	return r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO refresh_tokens(id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`,
			id, userID, token, expiresAt, createdAt)
		if err != nil {
			return fmt.Errorf("error creating refresh token: %w", err)
		}

		return nil
	})
}

func (r *RefreshTokenDBRepository) GetRefreshTokenByUserIDandToken(ctx context.Context, userID uuid.UUID, token string) (*domain.RefreshToken, error) {
	slog.Info("DB query", "operation", "SELECT", "table", "refresh_tokens", "parameters", fmt.Sprintf("user_id: %s, token: %s", userID, token))
	metrics.DBQueriesTotal.WithLabelValues("SELECT").Inc()

	var refreshToken domain.RefreshToken

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE user_id = $1 AND token = $2`,
		userID, token,
	).Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting refresh token by user ID and token: %w", err)
	}

	return &refreshToken, nil
}

func (r *RefreshTokenDBRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	slog.Info("DB query", "operation", "DELETE", "table", "refresh_tokens", "parameters", fmt.Sprintf("user_id: %s", userID))
	metrics.DBQueriesTotal.WithLabelValues("DELETE").Inc()

	return r.txProvider.Transact(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
		if err != nil {
			return fmt.Errorf("error revoking refresh tokens: %w", err)
		}

		return nil
	})
}
