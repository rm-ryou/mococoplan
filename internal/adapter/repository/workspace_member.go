package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type WorkspaceMemberRepository struct {
	db *sql.DB
}

func NewWorkspaceMemberRepository(db *sql.DB) ports.WorkspaceMemberRepository {
	return &WorkspaceMemberRepository{
		db: db,
	}
}

func (wmr *WorkspaceMemberRepository) Add(ctx context.Context, tx *sql.Tx, m *domain.WorkspaceMember) error {
	query := "INSERT INTO workspace_members (workspace_id, user_id, role) VALUES (?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, m.WorkspaceID, m.UserID, m.Role)
	return err
}

func (wmr *WorkspaceMemberRepository) Remove(ctx context.Context, tx *sql.Tx, workspaceID, userID int) error {
	query := "DELETE FROM workspace_members WHERE workspace_id = ? AND user_id = ?"

	_, err := tx.ExecContext(ctx, query, workspaceID, userID)
	return err
}

func (wmr *WorkspaceMemberRepository) FetchRole(ctx context.Context, workspaceID, userID int) (domain.WorkspaceRole, error) {
	query := "SELECT role FROM workspace_members WHERE workspace_id = ? AND user_id = ?"

	row := wmr.db.QueryRowContext(ctx, query, workspaceID, userID)

	var role string
	if err := row.Scan(&role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrNotFound
		}
		return "", err
	}
	return domain.WorkspaceRole(role), nil
}

func (wmr *WorkspaceMemberRepository) ListMembers(ctx context.Context, workspaceID int) ([]*domain.WorkspaceMember, error) {
	query := "SELECT * FROM workspace_members WHERE workspace_id = ?"

	rows, err := wmr.db.QueryContext(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.WorkspaceMember
	for rows.Next() {
		var m domain.WorkspaceMember
		var role string
		if err := rows.Scan(
			&m.WorkspaceID,
			&m.UserID,
			&role,
			&m.JoinedAt,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		m.Role = domain.WorkspaceRole(role)
		members = append(members, &m)
	}
	return members, nil
}

func (wmr *WorkspaceMemberRepository) Exists(ctx context.Context, workspaceID, userID int) (bool, error) {
	_, err := wmr.FetchRole(ctx, workspaceID, userID)
	return err == nil, err
}
