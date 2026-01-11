package ports

import (
	"context"
	"database/sql"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
)

type CreateChannelCmd struct {
	UserID      int
	WorkspaceID int
	Name        string
	Description string
	IsPrivate   bool
}

type ChannelServicer interface {
	Create(ctx context.Context, cmd *CreateChannelCmd) (*domain.Channel, error)
	ListJoined(ctx context.Context, workspaceID, userID int) ([]*domain.ChannelWithJoin, error)
	Join(ctx context.Context, channelID, userID int) error
	Leave(ctx context.Context, channelID, userID int) error
	AddMember(ctx context.Context, channelID, userID int) error
}

type ChannelRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ch *domain.Channel) (int, error)
	ListJoined(ctx context.Context, workspaceID, userID int) ([]*domain.ChannelWithJoin, error)
}

type ChannelMemberRepository interface {
	Add(ctx context.Context, tx *sql.Tx, m *domain.ChannelMember) error
	Remove(ctx context.Context, channelID int, userID int) error
	List(ctx context.Context, channelID int) ([]*domain.ChannelMember, error)
	Exists(ctx context.Context, channelID, userID int) (bool, error)
}
