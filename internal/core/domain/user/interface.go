package user

import "context"

type Creater interface {
	Create(ctx context.Context, u *User) error
}

type Servicer interface {
	Create(ctx context.Context, cmd *CreateCmd) error
}

type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type CreateCmd struct {
	Name          string
	Email         string
	PlainPassword string
}
