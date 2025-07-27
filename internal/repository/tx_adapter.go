package repository

import (
    "context"
    "TransactionTest/internal/domain"
)

type txAdapter struct {
    tx ITx
}

func (a *txAdapter) Commit(ctx context.Context) error {
    return a.tx.Commit(ctx)
}

func (a *txAdapter) Rollback(ctx context.Context) error {
    return a.tx.Rollback(ctx)
}

func (a *txAdapter) Exec(ctx context.Context, sql string, arguments ...interface{}) (domain.CommandTag, error) {
    return a.tx.Exec(ctx, sql, arguments...)
}

func (a *txAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) domain.Row {
    return rowAdapter{a.tx.QueryRow(ctx, sql, args...)}
}

type rowAdapter struct {
    row Row
}

func (r rowAdapter) Scan(dest ...interface{}) error {
    return r.row.Scan(dest...)
}