package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) ports.SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (sr *SessionRepository) Create(ctx context.Context, s *domain.Session) error {
	query := `
		INSERT INTO sessions
			(user_id, token, ip_address, user_agent, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := sr.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		s.UserID,
		s.Token[:],
		s.IP[:],
		s.UserAgent,
		s.ExpiresAt,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if ok := errors.As(err, &mysqlErr); ok {
			if mysqlErr.Number == 1062 {
				return domain.ErrSessionInvalid
			}
		}
		return err
	}

	return nil
}

func (sr *SessionRepository) FindByToken(ctx context.Context, token domain.SessionToken) (*domain.Session, error) {
	query := `
		SELECT
			id,
			user_id,
			token,
			ip_address,
			user_agent,
			expires_at,
			created_at,
			updated_at
		FROM
			sessions
		WHERE
			token = ?
	`

	row := sr.db.QueryRowContext(ctx, query, token[:])

	var tokenBytes []byte
	var ipBytes []byte
	var s domain.Session
	if err := row.Scan(
		&s.ID,
		&s.UserID,
		&tokenBytes,
		&ipBytes,
		&s.UserAgent,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}

	if len(tokenBytes) != 32 {
		return nil, domain.ErrSessionInvalid
	}
	copy(s.Token[:], tokenBytes)

	if len(ipBytes) != 16 {
		return nil, domain.ErrSessionInvalid
	}
	copy(s.IP[:], ipBytes)

	return &s, nil
}

func (sr *SessionRepository) Delete(ctx context.Context, token domain.SessionToken) error {
	query := "DELETE FROM sessions WHERE token = ?"

	_, err := sr.db.ExecContext(ctx, query, token[:])
	if err != nil {
		return err
	}

	return nil
}
