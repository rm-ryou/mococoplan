package auth

import "context"

type Servicer interface {
	Signup(ctx context.Context, cmd *SignupCmd) error
	Login(ctx context.Context, cmd *LoginCmd) error
	Logout(ctx context.Context) error
}

type LoginCmd struct {
	Email    string
	Password string
}

type SignupCmd struct {
	Name string
	LoginCmd
}
