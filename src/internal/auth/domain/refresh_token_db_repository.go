package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type RefreshTokenDBRepository interface {
	CreateRefreshToken(
		ctx context.Context,
		userID uuid.UUID,
		token string,
		expiresAt time.Time,
	) error
	GetRefreshTokenByUserIDandToken(ctx context.Context, userID uuid.UUID, token string) (*RefreshToken, error)
	RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}
