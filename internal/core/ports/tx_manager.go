package ports

import (
	"context"
	"database/sql"
)

type TxManager interface {
	WithinTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}
