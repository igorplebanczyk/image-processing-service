package auth

import "github.com/google/uuid"

type UserRepository interface {
	CreateUser(user *User) (uuid.UUID, error)
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByUsername(username string) (*User, error)
}
