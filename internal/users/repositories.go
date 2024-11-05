package users

import (
	"github.com/google/uuid"
	"time"
)

type UserRepository interface {
	GetUserByValue(field, value string) (*User, error)
	CreateUser(
		username string,
		email string,
		password string,
	) (*User, error)
}

type RefreshTokenRepository interface {
	GetRefreshTokenByValue(field, value string) (*RefreshToken, error)
	CreateRefreshToken(
		userID uuid.UUID,
		token string,
		expiresAt time.Time,
	) (*RefreshToken, error)
	RevokeRefreshToken(userID uuid.UUID) error
}
