package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type ChannelRepository struct {
	db *sql.DB
}

func NewChannelRepository(db *sql.DB) ports.ChannelRepository {
	return &ChannelRepository{
		db: db,
	}
}

func (cr *ChannelRepository) Create(ctx context.Context, tx *sql.Tx, ch *domain.Channel) (int, error) {
	query := "INSERT INTO channels (workspace_id, name, description, is_private) VALUES (?, ?, ?, ?)"

	res, err := tx.ExecContext(ctx, query, ch.WorkspaceID, ch.Name, ch.Description, ch.IsPrivate)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if ok := errors.As(err, &mysqlErr); ok {
			if mysqlErr.Number == 1062 {
				return 0, domain.ErrChannelAlreadyExists
			}
		}
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (cr *ChannelRepository) ListJoined(ctx context.Context, workspaceID, userID int) ([]*domain.ChannelWithJoin, error) {
	query := `
		SELECT
			c.id,
			c.name,
			c.description,
			c.is_private
		FROM
			channels c
		LEFT JOIN
			channel_members cm
		ON
			cm.channel_id = c.id AND cm.user_id = ?
		WHERE
			c.workspace_id = ?
		ORDER BY c.name
	`

	rows, err := cr.db.QueryContext(ctx, query, userID, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*domain.ChannelWithJoin
	for rows.Next() {
		var c domain.ChannelWithJoin
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.IsPrivate,
		); err != nil {
			return nil, err
		}
		channels = append(channels, &c)
	}

	return channels, nil
}
