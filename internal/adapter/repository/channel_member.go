package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type ChannelMemberRepository struct {
	db *sql.DB
}

func NewChannelMemberRepository(db *sql.DB) ports.ChannelMemberRepository {
	return &ChannelMemberRepository{
		db: db,
	}
}

func (cmr *ChannelMemberRepository) Add(ctx context.Context, tx *sql.Tx, m *domain.ChannelMember) error {
	query := "INSERT INTO channel_members (channel_id, user_id) VALUES (?, ?)"

	_, err := tx.ExecContext(ctx, query, m.ChannelID, m.UserID)
	return err
}

func (cmr *ChannelMemberRepository) Remove(ctx context.Context, channelID int, userID int) error {
	query := "DELETE FROM channel_members WHERE channel_id = ? AND user_id = ?"

	_, err := cmr.db.ExecContext(ctx, query, channelID, userID)
	return err
}

func (cmr *ChannelMemberRepository) List(ctx context.Context, channelID int) ([]*domain.ChannelMember, error) {
	query := "SELECT * FROM channel_members WHERE channel_id = ?"

	rows, err := cmr.db.QueryContext(ctx, query, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.ChannelMember
	for rows.Next() {
		var m domain.ChannelMember
		if err := rows.Scan(
			&m.ChannelID,
			&m.UserID,
			&m.JoinedAt,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	return members, nil
}

func (cmr *ChannelMemberRepository) Exists(ctx context.Context, channelID, userID int) (bool, error) {
	query := "SELECT 1 FROM channel_members WHERE channel_id = ? AND user_id = ?"

	row := cmr.db.QueryRowContext(ctx, query, channelID, userID)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, nil
	}
	return true, nil
}
