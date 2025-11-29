package port

import (
	"context"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
)

type UserServicer interface {
	Register(ctx context.Context, name, email, password string) (*domain.User, error)
}

type UserRepository interface {
	Register(ctx context.Context, user *domain.User) error
	List(ctx context.Context) ([]*domain.User, error)
	Delete(ctx context.Context, id domain.UserID) error
}
