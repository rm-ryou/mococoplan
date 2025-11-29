package user

import "context"

type Finder interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type Servicer interface {
	Create(ctx context.Context, name, email, password string) (*User, error)
}

type Repository interface {
	Create(ctx context.Context, u *User) (ID, error)
}
