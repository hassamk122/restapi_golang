package service

import "errors"

var (
	ErrEmailTaken            = errors.New("email already taken")
	ErrInvalidCredentials    = errors.New("Invalid credentials")
	ErrInvalidRequestPayload = errors.New("Invalid request payload")
	ErrUserNotFound          = errors.New("User not found")
)
