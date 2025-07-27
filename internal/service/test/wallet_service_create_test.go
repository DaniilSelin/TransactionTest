package test

import (
	"TransactionTest/internal/domain"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalletService_CreateWallet_Negative(t *testing.T) {
	ws := newWS(&MockWalletRepository{})
	addr, code := ws.CreateWallet(context.Background(), -1)
	assert.Equal(t, "", addr)
	assert.Equal(t, domain.CodeNegativeBalance, code)
}

func TestWalletService_CreateWallet_Duplicate(t *testing.T) {
	repo := &MockWalletRepository{
		CreateWalletFunc: func(ctx context.Context, address string, balance float64) error {
			return domain.ErrWalletAlreadyExists
		},
	}
	ws := newWS(repo)
	addr, code := ws.CreateWallet(context.Background(), 100)
	assert.Equal(t, "", addr)
	assert.Equal(t, domain.CodeDuplicateWallet, code)
}

func TestWalletService_CreateWallet_Internal(t *testing.T) {
	repo := &MockWalletRepository{
		CreateWalletFunc: func(ctx context.Context, address string, balance float64) error {
			return domain.ErrInternal
		},
	}
	ws := newWS(repo)
	addr, code := ws.CreateWallet(context.Background(), 100)
	assert.Equal(t, "", addr)
	assert.Equal(t, domain.CodeInternal, code)
}

func TestWalletService_CreateWallet_Success(t *testing.T) {
	repo := &MockWalletRepository{
		CreateWalletFunc: func(ctx context.Context, address string, balance float64) error {
			return nil
		},
	}
	ws := newWS(repo)
	addr, code := ws.CreateWallet(context.Background(), 100)
	assert.NotEmpty(t, addr)
	assert.Equal(t, domain.CodeOK, code)
}
