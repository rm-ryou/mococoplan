package domain

import (
	"time"
)

type User struct {
	ID            int
	Name          string
	Email         string
	EmailVerified bool
	PasswordHash  string
	ImageUrl      *string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
