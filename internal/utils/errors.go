package utils

import "errors"

var (
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("Invalid credentials")
)
