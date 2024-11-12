package domain

import "context"

type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}
