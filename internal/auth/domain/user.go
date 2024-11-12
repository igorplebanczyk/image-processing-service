package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
	Password string
}

func NewUser(iD uuid.UUID, username, password string) *User {
	return &User{
		ID:       iD,
		Username: username,
		Password: password,
	}
}
