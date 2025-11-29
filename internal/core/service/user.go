package service

import (
	"context"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/port"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) port.UserServicer {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	return nil, nil
}
