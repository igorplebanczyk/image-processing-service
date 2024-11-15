package domain

import "errors"

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInternal           = errors.New("internal server error")
)
