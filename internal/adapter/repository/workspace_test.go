package repository

import (
	"context"
	"testing"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func countWorkspacesByID(t *testing.T, id int) int {
	t.Helper()

	query := "SELECT COUNT(*) FROM workspaces WHERE id = ?"
	row := testDB.QueryRow(query, id)

	var cnt int
	require.NoError(t, row.Scan(&cnt))

	return cnt
}

func TestWorkspaceRepository_SuccessCreate(t *testing.T) {
	defer deleteAllRecords(t, testDB, "users")
	defer deleteAllRecords(t, testDB, "workspaces")
	userRepo := NewUserRepository(testDB)
	repo := NewWorkspaceRepository(testDB)

	err := userRepo.Create(context.Background(), &domain.User{
		Name:         "test name",
		Email:        "test@example.com",
		PasswordHash: "testHashedPassword",
	})
	require.NoError(t, err)

	user, err := userRepo.FindByEmail(context.Background(), "test@example.com")
	require.NoError(t, err)

	ws := &domain.Workspace{
		Name:      "Test",
		Slug:      "test-ws",
		CreatedBy: user.ID,
	}

	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	id, err := repo.Create(ctx, tx, ws)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	records := countWorkspacesByID(t, id)
	if records != 1 {
		t.Errorf("want workspaces count: 1, act: %d", records)
	}
}

func TestWorkspaceRepository_FailedCreate_DuplicateSlug(t *testing.T) {
	defer deleteAllRecords(t, testDB, "users")
	defer deleteAllRecords(t, testDB, "workspaces")
	userRepo := NewUserRepository(testDB)
	repo := NewWorkspaceRepository(testDB)

	err := userRepo.Create(context.Background(), &domain.User{
		Name:         "test name",
		Email:        "test@example.com",
		PasswordHash: "testHashedPassword",
	})
	require.NoError(t, err)

	user, err := userRepo.FindByEmail(context.Background(), "test@example.com")
	require.NoError(t, err)

	ws := &domain.Workspace{
		Name:      "Test",
		Slug:      "test-ws",
		CreatedBy: user.ID,
	}

	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	_, err = repo.Create(ctx, tx, ws)
	require.NoError(t, err)

	duplicatedSlug := &domain.Workspace{
		Name:      "Another",
		Slug:      "test-ws",
		CreatedBy: user.ID,
	}

	_, err = repo.Create(ctx, tx, duplicatedSlug)
	require.Error(t, err)
	if err != domain.ErrSlugAlreadyExists {
		t.Fatalf("want: %v, act: %v", domain.ErrSlugAlreadyExists, err)
	}
	require.NoError(t, tx.Commit())
}

func TestWorkspaceRepository_SuccessFindByID(t *testing.T) {
	defer deleteAllRecords(t, testDB, "users")
	defer deleteAllRecords(t, testDB, "workspaces")
	userRepo := NewUserRepository(testDB)
	repo := NewWorkspaceRepository(testDB)

	err := userRepo.Create(context.Background(), &domain.User{
		Name:         "test name",
		Email:        "test@example.com",
		PasswordHash: "testHashedPassword",
	})
	require.NoError(t, err)

	user, err := userRepo.FindByEmail(context.Background(), "test@example.com")
	require.NoError(t, err)

	ws := &domain.Workspace{
		Name:      "Test",
		Slug:      "test-ws",
		CreatedBy: user.ID,
	}

	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()

	id, err := repo.Create(ctx, tx, ws)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	act, err := repo.FindByID(ctx, id)
	require.NoError(t, err)

	if act.Slug != ws.Slug {
		t.Fatalf("want: %s, act: %s", ws.Slug, act.Slug)
	}
}

func TestWorkspaceRepository_FailedFindByID_NotFound(t *testing.T) {
	repo := NewWorkspaceRepository(testDB)

	act, err := repo.FindByID(context.Background(), -1)
	assert.Nil(t, act)
	assert.Equal(t, err, domain.ErrNotFound)
}
