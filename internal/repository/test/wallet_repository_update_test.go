package test

import (
	"context"
	"errors"
	"testing"
	"TransactionTest/internal/repository"
	"TransactionTest/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestWalletRepository_UpdateWalletBalance_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return &MockCommandTag{RowsAffectedFunc: func() int64 { return 1 }}, nil
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	err := repo.UpdateWalletBalance(ctx, "addr", 200)
	assert.NoError(t, err)
}

func TestWalletRepository_UpdateWalletBalance_NotFound(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return &MockCommandTag{RowsAffectedFunc: func() int64 { return 0 }}, nil
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	err := repo.UpdateWalletBalance(ctx, "addr", 200)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestWalletRepository_UpdateWalletBalance_NegativeBalance(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return nil, &mockDBError{sqlState: repository.ErrCodeCheckViolation, constraint: repository.ConstraintBalanceNonNegative}
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	err := repo.UpdateWalletBalance(ctx, "addr", -100)
	assert.True(t, errors.Is(err, domain.ErrNegativeBalance))
}

func TestWalletRepository_UpdateWalletBalance_InternalError(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return nil, errors.New("fail")
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	err := repo.UpdateWalletBalance(ctx, "addr", 200)
	assert.True(t, errors.Is(err, domain.ErrInternal))
}

func TestWalletRepository_UpdateWalletBalanceTx_Success(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return &MockCommandTag{RowsAffectedFunc: func() int64 { return 1 }}, nil
		},
	}
	repo := repository.NewWalletRepository(nil)
	err := repo.UpdateWalletBalanceTx(ctx, *mockTx, "addr", 200)
	assert.NoError(t, err)
}

func TestWalletRepository_UpdateWalletBalanceTx_NotFound(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return &MockCommandTag{RowsAffectedFunc: func() int64 { return 0 }}, nil
		},
	}
	repo := repository.NewWalletRepository(nil)
	err := repo.UpdateWalletBalanceTx(ctx, *mockTx, "addr", 200)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestWalletRepository_UpdateWalletBalanceTx_NegativeBalance(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return nil, &mockDBError{sqlState: repository.ErrCodeCheckViolation, constraint: repository.ConstraintBalanceNonNegative}
		},
	}
	repo := repository.NewWalletRepository(nil)
	err := repo.UpdateWalletBalanceTx(ctx, *mockTx, "addr", -100)
	assert.True(t, errors.Is(err, domain.ErrNegativeBalance))
}

func TestWalletRepository_UpdateWalletBalanceTx_InternalError(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return nil, errors.New("fail")
		},
	}
	repo := repository.NewWalletRepository(nil)
	err := repo.UpdateWalletBalanceTx(ctx, *mockTx, "addr", 200)
	assert.True(t, errors.Is(err, domain.ErrInternal))
}