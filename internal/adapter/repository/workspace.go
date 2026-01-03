package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type WorkspaceRepository struct {
	db *sql.DB
}

func NewWorkspaceRepository(db *sql.DB) ports.WorkspaceRepository {
	return &WorkspaceRepository{
		db: db,
	}
}

func (wr *WorkspaceRepository) Create(ctx context.Context, tx *sql.Tx, ws *domain.Workspace) (int, error) {
	query := "INSERT INTO workspaces (name, slug, created_by) VALUES (?, ?, ?)"

	res, err := tx.ExecContext(ctx, query, ws.Name, ws.Slug, ws.CreatedBy)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if ok := errors.As(err, &mysqlErr); ok {
			if mysqlErr.Number == 1062 {
				return 0, domain.ErrSlugAlreadyExists
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

func (wr *WorkspaceRepository) FindByID(ctx context.Context, workspaceID int) (*domain.Workspace, error) {
	query := "SELECT * FROM workspaces WHERE id = ?"

	row := wr.db.QueryRowContext(ctx, query, workspaceID)

	var ws domain.Workspace
	if err := row.Scan(
		&ws.ID,
		&ws.Name,
		&ws.Slug,
		&ws.CreatedBy,
		&ws.CreatedAt,
		&ws.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &ws, nil
}
func (wr *WorkspaceRepository) ListByUser(ctx context.Context, userID int) ([]*domain.Workspace, error) {
	query := `
		SELECT
			ws.id,
			ws.name,
			ws.slug,
			ws.created_by,
			ws.created_at,
			ws.updated_at
		FROM
			workspaces ws
		JOIN
			workspace_members wsm
		ON
			wsm.workspace_id = ws.id
		WHERE
			wsm.user_id = ?
		ORDER BY
			ws.created_at DESC
	`

	rows, err := wr.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*domain.Workspace
	for rows.Next() {
		var ws domain.Workspace
		if err := rows.Scan(
			&ws.ID,
			&ws.Name,
			&ws.Slug,
			&ws.CreatedBy,
			&ws.CreatedAt,
			&ws.UpdatedAt,
		); err != nil {
			return nil, err
		}
		fmt.Println("ws:", ws)
		workspaces = append(workspaces, &ws)
	}

	return workspaces, nil
}
