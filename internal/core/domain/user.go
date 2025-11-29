package domain

import (
	"time"
)

type User struct {
	Id            int
	Name          string
	Email         string
	EmailVerified bool
	ImageUrl      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
