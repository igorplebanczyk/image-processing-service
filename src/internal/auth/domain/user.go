package domain

import (
	"github.com/google/uuid"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	ID       uuid.UUID
	Username string
	Password string
	Role     Role
}
