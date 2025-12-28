package ports

import (
	"context"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
