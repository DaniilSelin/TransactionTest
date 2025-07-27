package repository

import "context.Context"

type txAdapter struct {
    tx ITx
}

func (a *txAdapter) Commit(ctx context.Context) error {
    return a.tx.Commit(ctx)
}

func (a *txAdapter) Rollback(ctx context.Context) error {
    return a.tx.Rollback(ctx)
}