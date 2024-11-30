package domain

import (
	"context"
	"github.com/google/uuid"
)

type UserDBRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserRoleByID(ctx context.Context, id uuid.UUID) (Role, error)
}
