package users

import (
	"github.com/google/uuid"
	"time"
)

type UserRepository interface {
	CreateUser(username, email, password string) (*User, error)
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(id uuid.UUID, username, email string) error
	DeleteUser(id uuid.UUID) error
}

type RefreshTokenRepository interface {
	CreateRefreshToken(
		userID uuid.UUID,
		token string,
		expiresAt time.Time,
	) (*RefreshToken, error)
	GetRefreshTokenByUserID(userID uuid.UUID) (*RefreshToken, error)
	RevokeRefreshToken(userID uuid.UUID) error
}
