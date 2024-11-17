package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(
		ctx context.Context,
		userID uuid.UUID,
		token string,
		expiresAt time.Time,
	) error
	GetRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) ([]*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, userID uuid.UUID) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}
