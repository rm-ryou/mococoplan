package service

import (
	"context"

	"github.com/rm-ryou/mococoplan/internal/core/domain/user"
	"github.com/rm-ryou/mococoplan/pkg/password"
)

type UserService struct {
	repo   user.Repository
	params *password.Params
}

func NewUserService(repo user.Repository, params *password.Params) user.Servicer {
	return &UserService{
		repo:   repo,
		params: params,
	}
}

func (us *UserService) Create(ctx context.Context, cmd *user.CreateCmd) error {
	hash, err := password.Hash(cmd.PlainPassword, us.params)
	if err != nil {
		return err
	}

	u := &user.User{
		Name:          cmd.Name,
		Email:         cmd.Email,
		EmailVerified: false,
		PasswordHash:  hash,
	}

	if err := us.repo.Create(ctx, u); err != nil {
		return err
	}

	return nil
}
