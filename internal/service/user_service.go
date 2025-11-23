package service

import (
	"context"

	"github.com/rm-ryou/mococoplan/internal/api/controller"
	"github.com/rm-ryou/mococoplan/internal/domain"
)

type UserRepository interface {
	Register(ctx context.Context, user *domain.User) error
	List(ctx context.Context) ([]*domain.User, error)
	Delete(ctx context.Context, id domain.UserID) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) controller.UserServicer {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	return nil, nil
}
