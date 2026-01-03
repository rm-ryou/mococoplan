package repository

import (
	"context"
	"database/sql"

	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) ports.TxManager {
	return &TxManager{
		db: db,
	}
}

func (tm *TxManager) WithinTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
