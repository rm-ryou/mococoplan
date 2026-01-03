package service

import (
	"context"
	"database/sql"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type WorkspaceService struct {
	tx      ports.TxManager
	wsRepo  ports.WorkspaceRepository
	wsmRepo ports.WorkspaceMemberRepository
}

func NewWorkspaceService(
	tx ports.TxManager,
	wsRepo ports.WorkspaceRepository,
	wsmRepo ports.WorkspaceMemberRepository,
) ports.WorkspaceServicer {
	return &WorkspaceService{
		tx:      tx,
		wsRepo:  wsRepo,
		wsmRepo: wsmRepo,
	}
}

func (ws *WorkspaceService) Create(ctx context.Context, cmd *ports.CreateWorkspaceCmd) (*domain.Workspace, error) {
	workspace := &domain.Workspace{
		Name:      cmd.Name,
		Slug:      cmd.Slug,
		CreatedBy: cmd.UserID,
	}

	if err := ws.tx.WithinTx(ctx, func(tx *sql.Tx) error {
		workspaceID, err := ws.wsRepo.Create(ctx, tx, workspace)
		if err != nil {
			return err
		}
		workspace.ID = workspaceID

		owner := &domain.WorkspaceMember{
			WorkspaceID: workspace.ID,
			UserID:      workspace.CreatedBy,
			Role:        domain.WorkspaceRoleOwner,
		}

		if err := ws.wsmRepo.Add(ctx, tx, owner); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (ws *WorkspaceService) AddMember(ctx context.Context, cmd *ports.AddWorkspaceMemberCmd) error {
	role, err := ws.wsmRepo.FetchRole(ctx, cmd.WorkspaceID, cmd.UserID)
	if err != nil {
		return nil
	}
	if role != domain.WorkspaceRoleOwner && role != domain.WorkspaceRoleAdmin {
		return domain.ErrForbiddenRole
	}

	exists, err := ws.wsmRepo.Exists(ctx, cmd.WorkspaceID, cmd.TargetUserID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	member := &domain.WorkspaceMember{
		WorkspaceID: cmd.WorkspaceID,
		UserID:      cmd.TargetUserID,
		Role:        cmd.Role,
	}

	return ws.tx.WithinTx(ctx, func(tx *sql.Tx) error {
		return ws.wsmRepo.Add(ctx, tx, member)
	})
}

func (ws *WorkspaceService) ListWorkspaces(ctx context.Context, userID int) ([]*domain.Workspace, error) {
	return ws.wsRepo.ListByUser(ctx, userID)
}

func (ws *WorkspaceService) ListMembers(ctx context.Context, workspaceID, userID int) ([]*domain.WorkspaceMember, error) {
	return ws.wsmRepo.ListMembers(ctx, workspaceID)
}
