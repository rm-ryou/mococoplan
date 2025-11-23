package domain

import (
	"time"
)

type UserID int

type User struct {
	ID            UserID
	Name          string
	Email         string
	EmailVerified bool
	ImageUrl      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
