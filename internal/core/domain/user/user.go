package user

import "time"

type ID int

type User struct {
	Id            ID
	Name          string
	Email         string
	EmailVerified bool
	PasswordHash  string
	ImageUrl      *string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
