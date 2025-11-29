package user

import "time"

type User struct {
	Id            int
	Name          string
	Email         string
	EmailVerified bool
	PasswordHash  string
	ImageUrl      *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
