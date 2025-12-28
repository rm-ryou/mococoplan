package ports

import (
	"context"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, s *domain.Session) error
	FindByToken(ctx context.Context, token domain.SessionToken) (*domain.Session, error)
	Delete(ctx context.Context, token domain.SessionToken) error
}
