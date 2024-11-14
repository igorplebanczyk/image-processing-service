package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInternal           = errors.New("internal server error")
)
