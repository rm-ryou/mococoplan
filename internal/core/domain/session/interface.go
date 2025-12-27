package session

import "context"

type Repository interface {
	Create(ctx context.Context, s *Session) error
	FindByToken(ctx context.Context, token [32]byte) (*Session, error)
	Delete(ctx context.Context, token [32]byte) error
}
