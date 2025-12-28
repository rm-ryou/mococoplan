package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := "INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)"

	stmt, err := ur.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		u.Name,
		u.Email,
		u.PasswordHash,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if ok := errors.As(err, &mysqlErr); ok {
			if mysqlErr.Number == 1062 {
				return domain.ErrEmailAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT
			id,
			name,
			email,
			email_verified,
			password_hash,
			image_url,
			created_at,
			updated_at
		FROM
			users
		WHERE
			email = ?
	`
	row := ur.db.QueryRowContext(ctx, query, email)

	var u domain.User
	if err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.EmailVerified,
		&u.PasswordHash,
		&u.ImageUrl,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}
