package token

import (
	"time"
)

type Issuer interface {
	Issue(claims *Claims) (*AccessToken, error)
}

type Verifier interface {
	Verify(token string) (*Claims, error)
}

type Servicer interface {
	Issuer
	Verifier
}

type Claims struct {
	UserId int
	Email  string
}

type AccessToken struct {
	Token     string
	ExpiresAt time.Time
}
