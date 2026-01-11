package service

import (
	"context"
	"database/sql"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type ChannelService struct {
	tx      ports.TxManager
	chRepo  ports.ChannelRepository
	chmRepo ports.ChannelMemberRepository
}

func NewChannelService(
	tx ports.TxManager,
	chRepo ports.ChannelRepository,
	chmRepo ports.ChannelMemberRepository,
) ports.ChannelServicer {
	return &ChannelService{
		tx:      tx,
		chRepo:  chRepo,
		chmRepo: chmRepo,
	}
}

func (cs *ChannelService) Create(ctx context.Context, cmd *ports.CreateChannelCmd) (*domain.Channel, error) {
	ch := &domain.Channel{
		WorkspaceID: cmd.WorkspaceID,
		Name:        cmd.Name,
		Description: cmd.Description,
		IsPrivate:   cmd.IsPrivate,
		CreatedBy:   cmd.UserID,
	}

	if err := cs.tx.WithinTx(ctx, func(tx *sql.Tx) error {
		channelID, err := cs.chRepo.Create(ctx, tx, ch)
		if err != nil {
			return err
		}
		ch.ID = channelID

		chMember := &domain.ChannelMember{
			ChannelID: channelID,
			UserID:    cmd.UserID,
		}

		return cs.chmRepo.Add(ctx, tx, chMember)
	}); err != nil {
		return nil, err
	}

	return ch, nil
}

func (cs *ChannelService) ListJoined(ctx context.Context, workspaceID, userID int) ([]*domain.ChannelWithJoin, error) {
	return cs.chRepo.ListJoined(ctx, workspaceID, userID)
}

func (cs *ChannelService) Join(ctx context.Context, channelID, userID int) error {
	ok, err := cs.chmRepo.Exists(ctx, channelID, userID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	member := &domain.ChannelMember{
		ChannelID: channelID,
		UserID:    userID,
	}

	return cs.tx.WithinTx(ctx, func(tx *sql.Tx) error {
		return cs.chmRepo.Add(ctx, tx, member)
	})
}

func (cs *ChannelService) Leave(ctx context.Context, channelID, userID int) error {
	return cs.chmRepo.Remove(ctx, channelID, userID)
}

func (cs *ChannelService) AddMember(ctx context.Context, channelID, userID int) error {
	member := &domain.ChannelMember{
		ChannelID: channelID,
		UserID:    userID,
	}

	return cs.tx.WithinTx(ctx, func(tx *sql.Tx) error {
		return cs.chmRepo.Add(ctx, tx, member)
	})
}
