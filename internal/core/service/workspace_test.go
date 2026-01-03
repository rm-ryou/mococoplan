package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
	"github.com/rm-ryou/mococoplan/internal/core/ports/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceService_SuccessCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := new(mocks.TxManager)
	wsRepo := new(mocks.WorkspaceRepository)
	wsmRepo := new(mocks.WorkspaceMemberRepository)

	service := NewWorkspaceService(tx, wsRepo, wsmRepo)

	cmd := &ports.CreateWorkspaceCmd{
		UserID: 1,
		Name:   "test-name",
		Slug:   "test-slug",
	}

	tx.On("WithinTx", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(*sql.Tx) error)
			_ = fn(nil)
		}).
		Return(nil).Once()

	workspace := &domain.Workspace{
		Name:      cmd.Name,
		Slug:      cmd.Slug,
		CreatedBy: cmd.UserID,
	}
	wsRepo.On("Create", ctx, mock.Anything, workspace).Return(1, nil).Once()

	owner := &domain.WorkspaceMember{
		WorkspaceID: 1,
		UserID:      cmd.UserID,
		Role:        domain.WorkspaceRoleOwner,
	}
	wsmRepo.On("Add", ctx, mock.Anything, owner).Return(nil).Once()

	ws, err := service.Create(ctx, cmd)
	require.NoError(t, err)
	require.NotNil(t, ws)

	tx.AssertExpectations(t)
	wsRepo.AssertExpectations(t)
	wsmRepo.AssertExpectations(t)
}

func TestWorkspaceService_SuccessAddMember(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := new(mocks.TxManager)
	wsRepo := new(mocks.WorkspaceRepository)
	wsmRepo := new(mocks.WorkspaceMemberRepository)

	service := NewWorkspaceService(tx, wsRepo, wsmRepo)

	cmd := &ports.AddWorkspaceMemberCmd{
		WorkspaceID:  1,
		UserID:       1,
		TargetUserID: 2,
		Role:         domain.WorkspaceRoleMember,
	}

	tx.On("WithinTx", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(*sql.Tx) error)
			_ = fn(nil)
		}).
		Return(nil).Once()

	wsmRepo.On("FetchRole", ctx, cmd.WorkspaceID, cmd.UserID).Return(domain.WorkspaceRoleOwner, nil).Once()
	wsmRepo.On("Exists", ctx, cmd.WorkspaceID, cmd.TargetUserID).Return(false, nil).Once()

	member := &domain.WorkspaceMember{
		WorkspaceID: cmd.WorkspaceID,
		UserID:      cmd.TargetUserID,
		Role:        cmd.Role,
	}
	wsmRepo.On("Add", ctx, mock.Anything, member).Return(nil).Once()

	err := service.AddMember(ctx, cmd)
	require.NoError(t, err)

	tx.AssertExpectations(t)
	wsRepo.AssertExpectations(t)
	wsmRepo.AssertExpectations(t)
}

func TestWorkspaceService_FailedAddMember_WhenInvalidRole(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := new(mocks.TxManager)
	wsRepo := new(mocks.WorkspaceRepository)
	wsmRepo := new(mocks.WorkspaceMemberRepository)

	service := NewWorkspaceService(tx, wsRepo, wsmRepo)

	cmd := &ports.AddWorkspaceMemberCmd{
		WorkspaceID:  1,
		UserID:       1,
		TargetUserID: 2,
		Role:         domain.WorkspaceRoleMember,
	}

	wsmRepo.On("FetchRole", ctx, cmd.WorkspaceID, cmd.UserID).Return(domain.WorkspaceRoleMember, nil).Once()

	err := service.AddMember(ctx, cmd)
	require.Error(t, err)

	assert.ErrorIs(t, err, domain.ErrForbiddenRole)

	tx.AssertNotCalled(t, "WithinTx", ctx, mock.Anything)
	wsmRepo.AssertExpectations(t)
	wsmRepo.AssertNotCalled(t, "Exists", ctx, cmd.WorkspaceID, cmd.TargetUserID)
	wsmRepo.AssertNotCalled(t, "Add", ctx, mock.Anything, mock.Anything)
}
