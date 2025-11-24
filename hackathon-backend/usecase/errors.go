package usecase

import "errors"

var (
	ErrInvalidUser = errors.New("invalid user")
	ErrEmailExists = errors.New("email already exists")
)
