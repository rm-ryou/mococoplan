package domain

import "errors"

var (
	ErrNotFound = errors.New("not found")

	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionRevoked  = errors.New("session revoked")
	ErrSessionInvalid  = errors.New("invalid session")

	ErrSlugAlreadyExists = errors.New("slug already exists")
	ErrForbiddenRole     = errors.New("forbidden role")
)
