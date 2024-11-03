package auth

import "github.com/google/uuid"

type UserRepository interface {
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(
		username string,
		email string,
		password string,
	) (*User, error)
}
