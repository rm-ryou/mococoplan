package ports

import (
	"context"
	"database/sql"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
)

type CreateWorkspaceCmd struct {
	UserID int
	Name   string
	Slug   string
}

type AddWorkspaceMemberCmd struct {
	UserID       int
	WorkspaceID  int
	TargetUserID int
	Role         domain.WorkspaceRole
}

type WorkspaceServicer interface {
	Create(ctx context.Context, cmd *CreateWorkspaceCmd) (*domain.Workspace, error)
	AddMember(ctx context.Context, cmd *AddWorkspaceMemberCmd) error
	ListWorkspaces(ctx context.Context, userID int) ([]*domain.Workspace, error)
	ListMembers(ctx context.Context, workspaceID, userID int) ([]*domain.WorkspaceMember, error)
}

type WorkspaceRepository interface {
	Create(ctx context.Context, tx *sql.Tx, ws *domain.Workspace) (int, error)
	FindByID(ctx context.Context, workspaceID int) (*domain.Workspace, error)
	ListByUser(ctx context.Context, userID int) ([]*domain.Workspace, error)
}

type WorkspaceMemberRepository interface {
	Add(ctx context.Context, tx *sql.Tx, m *domain.WorkspaceMember) error
	Remove(ctx context.Context, tx *sql.Tx, workspaceID, userID int) error
	FetchRole(ctc context.Context, workspaceID, userID int) (domain.WorkspaceRole, error)
	ListMembers(ctx context.Context, workspaceID int) ([]*domain.WorkspaceMember, error)
	Exists(ctx context.Context, workspaceID, userID int) (bool, error)
}
