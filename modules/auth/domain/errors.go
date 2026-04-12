package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrDuplicateKey    = errors.New("duplicate key")
	ErrInvalidEntity   = errors.New("invalid entity")
	ErrInvalidID       = errors.New("invalid id")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInvalidPassword = errors.New("invalid password")
	ErrEmailExists     = errors.New("email already exists")
	ErrInactiveUser    = errors.New("inactive user")
)
