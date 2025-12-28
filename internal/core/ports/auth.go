package ports

import (
	"context"
	"net"
	"time"
)

type LoginCmd struct {
	Email     string
	Password  string
	IP        net.IP
	UserAgent string
}

type LogoutCmd struct {
	Token string
}

type SignupCmd struct {
	Name string
	LoginCmd
}

type AuthResult struct {
	AccessToken  string
	ExpiresAt    time.Time
	RefreshToken string
}

type AuthServicer interface {
	Login(ctx context.Context, cmd *LoginCmd) (*AuthResult, error)
	Logout(ctx context.Context, cmd *LogoutCmd) error
	Signup(ctx context.Context, cmd *SignupCmd) (*AuthResult, error)
	RefreshTokenTTL() time.Duration
}
