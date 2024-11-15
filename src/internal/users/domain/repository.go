package domain

import (
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(
		ctx context.Context,
		username,
		email,
		password string,
	) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUser(
		ctx context.Context,
		id uuid.UUID,
		username,
		email string,
	) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
