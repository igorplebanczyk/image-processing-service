package domain

import (
	"context"
	"github.com/google/uuid"
)

type UsersDBRepository interface {
	CreateUser(
		ctx context.Context,
		username,
		email,
		password,
		otpSecret string,
	) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetAllUsers(ctx context.Context, page, limit int) ([]User, int, error)
	UpdateUserDetails(
		ctx context.Context,
		id uuid.UUID,
		username,
		email string,
	) error
	UpdateUserRole(ctx context.Context, id uuid.UUID, role Role) error
	UpdateUserAsVerified(ctx context.Context, id uuid.UUID) error
	UpdateUserPassword(ctx context.Context, id uuid.UUID, password string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
