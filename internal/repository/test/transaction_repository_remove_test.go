package test

import (
	"TransactionTest/internal/domain"
	"TransactionTest/internal/repository"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransactionRepository_RemoveTransaction_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return &MockCommandTag{RowsAffectedFunc: func() int64 { return 1 }}, nil
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	err := repo.RemoveTransaction(ctx, 1)
	assert.NoError(t, err)
}

func TestTransactionRepository_RemoveTransaction_NotFound(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return &MockCommandTag{RowsAffectedFunc: func() int64 { return 0 }}, nil
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	err := repo.RemoveTransaction(ctx, 1)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestTransactionRepository_RemoveTransaction_InternalError(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return nil, errors.New("fail")
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	err := repo.RemoveTransaction(ctx, 1)
	assert.True(t, errors.Is(err, domain.ErrInternal))
}
