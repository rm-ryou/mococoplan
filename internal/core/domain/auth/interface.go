package auth

import "context"

type Servicer interface {
	Login(ctx context.Context, email, password string) error
	Logout(ctx context.Context) error
}
