package domain

import (
	"context.Context"
)

type TxExecutor interface {
    Commit(ctx context.Context) error
    Rollback(ctx context.Context) error
}