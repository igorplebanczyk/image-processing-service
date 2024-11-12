package domain

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func NewRefreshToken(iD, userID uuid.UUID, token string, expiresAt, createdAt time.Time) *RefreshToken {
	return &RefreshToken{
		ID:        iD,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}
}
