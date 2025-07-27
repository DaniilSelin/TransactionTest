package test

import (
	"context"
	"errors"
	"testing"
	"TransactionTest/internal/repository"
	"TransactionTest/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestWalletRepository_GetWalletBalance_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				*dest[0].(*float64) = 123.45
				return nil
			}}
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	balance, err := repo.GetWalletBalance(ctx, "addr")
	assert.NoError(t, err)
	assert.Equal(t, 123.45, balance)
}

func TestWalletRepository_GetWalletBalance_NotFound(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return &mockDBError{sqlState: "no_rows"}
			}}
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	_, err := repo.GetWalletBalance(ctx, "addr")
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestWalletRepository_GetWalletBalance_InternalError(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return errors.New("fail")
			}}
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	_, err := repo.GetWalletBalance(ctx, "addr")
	assert.True(t, errors.Is(err, domain.ErrInternal))
}

func TestWalletRepository_GetWallet_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				*dest[0].(*string) = "addr"
				*dest[1].(*float64) = 123.45
				return nil
			}}
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	wallet, err := repo.GetWallet(ctx, "addr")
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, "addr", wallet.Address)
	assert.Equal(t, 123.45, wallet.Balance)
}

func TestWalletRepository_GetWallet_NotFound(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return &mockDBError{sqlState: "no_rows"}
			}}
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	wallet, err := repo.GetWallet(ctx, "addr")
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	assert.Nil(t, wallet)
}

func TestWalletRepository_GetWallet_InternalError(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return errors.New("fail")
			}}
		},
	}
	repo := repository.NewWalletRepository(mockDB)
	wallet, err := repo.GetWallet(ctx, "addr")
	assert.True(t, errors.Is(err, domain.ErrInternal))
	assert.Nil(t, wallet)
}