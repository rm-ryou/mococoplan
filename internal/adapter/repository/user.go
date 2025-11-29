package repository

import (
	"context"
	"database/sql"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/port"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) port.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Register(ctx context.Context, user *domain.User) error {
	return nil
}

func (ur *UserRepository) List(ctx context.Context) ([]*domain.User, error) {
	return nil, nil
}

func (ur *UserRepository) Delete(ctx context.Context, id domain.UserID) error {
	return nil
}
