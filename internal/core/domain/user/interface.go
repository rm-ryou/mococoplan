package user

import "context"

type Creater interface {
	Create(ctx context.Context, u *User) error
}

type Servicer interface {
	Creater
}

type Repository interface {
	Creater
	FindByEmail(ctx context.Context, email string) (*User, error)
}
