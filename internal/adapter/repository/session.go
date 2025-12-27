package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/rm-ryou/mococoplan/internal/core/domain/session"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) session.Repository {
	return &SessionRepository{
		db: db,
	}
}

func (sr *SessionRepository) Create(ctx context.Context, s *session.Session) error {
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
		s.UserId,
		s.Token[:],
		s.IPAddress[:],
		s.UserAgent,
		s.ExpiresAt,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if ok := errors.As(err, &mysqlErr); ok {
			if mysqlErr.Number == 1062 {
				return session.ErrInvalid
			}
		}
		return err
	}

	return nil
}

func (sr *SessionRepository) FindByToken(ctx context.Context, token [32]byte) (*session.Session, error) {
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
	var ipAddressBytes []byte
	var s session.Session
	if err := row.Scan(
		&s.Id,
		&s.UserId,
		&tokenBytes,
		&ipAddressBytes,
		&s.UserAgent,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, session.ErrNotFound
		}
		return nil, err
	}

	if len(tokenBytes) != 32 {
		return nil, session.ErrInvalid
	}
	copy(s.Token[:], tokenBytes)

	if len(ipAddressBytes) != 16 {
		return nil, session.ErrInvalid
	}
	copy(s.IPAddress[:], ipAddressBytes)

	return &s, nil
}

func (sr *SessionRepository) Delete(ctx context.Context, token [32]byte) error {
	query := "DELETE FROM sessions WHERE token = ?"

	_, err := sr.db.ExecContext(ctx, query, token[:])
	if err != nil {
		return err
	}

	return nil
}
