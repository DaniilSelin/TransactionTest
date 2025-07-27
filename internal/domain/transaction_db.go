package domain

import (
	"context"
)

type TxExecutor interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row
}

type Row interface {
	Scan(dest ...interface{}) error
}

type CommandTag interface {
	RowsAffected() int64
}
