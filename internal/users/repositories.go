package users

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type UserRepository interface {
	CreateUser(
		ctx context.Context,
		username,
		email,
		password string,
	) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(
		ctx context.Context,
		id uuid.UUID,
		username,
		email string,
	) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type RefreshTokenRepository interface {
	CreateRefreshToken(
		ctx context.Context,
		userID uuid.UUID,
		token string,
		expiresAt time.Time,
	) (*RefreshToken, error)
	GetRefreshTokenByUserID(ctx context.Context, userID uuid.UUID) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, userID uuid.UUID) error
}
