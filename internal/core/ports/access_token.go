package ports

import "time"

type UserIdentity struct {
	UserID int
	Email  string
}

type AccessToken struct {
	Token     string
	ExpiresAt time.Time
}

type TokenIssuer interface {
	Issue(identity *UserIdentity) (*AccessToken, error)
}

type TokenVerifier interface {
	Verify(token string) (*UserIdentity, error)
}
