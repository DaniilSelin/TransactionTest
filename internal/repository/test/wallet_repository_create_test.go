package test

import (
	"TransactionTest/internal/domain"
	"TransactionTest/internal/repository"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalletRepository_CreateWallet_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return &MockCommandTag{RowsAffectedFunc: func() int64 { return 1 }}, nil
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	err := repo.CreateWallet(ctx, "addr", 100)
	assert.NoError(t, err)
}

func TestWalletRepository_CreateWallet_Duplicate(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return nil, &mockDBError{sqlState: repository.ErrCodeUniqueViolation}
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	err := repo.CreateWallet(ctx, "addr", 100)
	assert.True(t, errors.Is(err, domain.ErrWalletAlreadyExists))
}

func TestWalletRepository_CreateWallet_NegativeBalance(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return nil, &mockDBError{sqlState: repository.ErrCodeCheckViolation, constraint: repository.ConstraintBalanceNonNegative}
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	err := repo.CreateWallet(ctx, "addr", -100)
	assert.True(t, errors.Is(err, domain.ErrNegativeBalance))
}

func TestWalletRepository_CreateWallet_InternalError(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return nil, errors.New("db fail")
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	err := repo.CreateWallet(ctx, "addr", 100)
	assert.True(t, errors.Is(err, domain.ErrInternal))
}

func TestWalletRepository_CreateWalletTx_Success(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return &MockCommandTag{RowsAffectedFunc: func() int64 { return 1 }}, nil
		},
	}
	repo := repository.NewWalletRepository(nil)
	err := repo.CreateWalletTx(ctx, mockTx, "addr", 100)
	assert.NoError(t, err)
}

func TestWalletRepository_CreateWalletTx_Duplicate(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return nil, &mockDBError{sqlState: repository.ErrCodeUniqueViolation}
		},
	}
	repo := repository.NewWalletRepository(nil)
	err := repo.CreateWalletTx(ctx, mockTx, "addr", 100)
	assert.True(t, errors.Is(err, domain.ErrWalletAlreadyExists))
}

func TestWalletRepository_CreateWalletTx_NegativeBalance(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return nil, &mockDBError{sqlState: repository.ErrCodeCheckViolation, constraint: repository.ConstraintBalanceNonNegative}
		},
	}
	repo := repository.NewWalletRepository(nil)
	err := repo.CreateWalletTx(ctx, mockTx, "addr", -100)
	assert.True(t, errors.Is(err, domain.ErrNegativeBalance))
}

func TestWalletRepository_CreateWalletTx_InternalError(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.CommandTag, error) {
			return nil, errors.New("fail")
		},
	}
	repo := repository.NewWalletRepository(nil)
	err := repo.CreateWalletTx(ctx, mockTx, "addr", 100)
	assert.True(t, errors.Is(err, domain.ErrInternal))
}
