package domain

import "errors"

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrValidationFailed = errors.New("validation failed")
	ErrInternal         = errors.New("internal server error")
)
