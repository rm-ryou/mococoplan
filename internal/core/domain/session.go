package domain

import (
	"time"
)

type SessionToken [32]byte
type SessionIP [16]byte

type Session struct {
	ID        int
	UserID    int
	Token     SessionToken
	IP        SessionIP
	UserAgent string
	ExpiresAt time.Time
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (s *Session) IsExpired(now time.Time) bool {
	return !s.ExpiresAt.After(now)
}
