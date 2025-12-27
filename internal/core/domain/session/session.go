package session

import (
	"errors"
	"time"
)

type Session struct {
	Id        int
	UserId    int
	Token     [32]byte
	IPAddress [16]byte
	UserAgent string
	ExpiresAt time.Time
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

var (
	ErrNotFound = errors.New("session not found")
	ErrExpired  = errors.New("session expired")
	ErrRevoked  = errors.New("session revoked")
	ErrInvalid  = errors.New("invalid session")
)

func (s *Session) IsExpired(now time.Time) bool {
	return !s.ExpiresAt.After(now)
}
