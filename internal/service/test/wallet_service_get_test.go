package test

import (
	"TransactionTest/internal/domain"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalletService_GetBalance_NotFound(t *testing.T) {
	repo := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			return 0, domain.ErrNotFound
		},
	}
	ws := newWS(repo)
	bal, code := ws.GetBalance(context.Background(), "addr")
	assert.Equal(t, float64(0), bal)
	assert.Equal(t, domain.CodeWalletNotFound, code)
}

func TestWalletService_GetBalance_Internal(t *testing.T) {
	repo := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			return 0, errors.New("fail")
		},
	}
	ws := newWS(repo)
	bal, code := ws.GetBalance(context.Background(), "addr")
	assert.Equal(t, float64(0), bal)
	assert.Equal(t, domain.CodeInternal, code)
}

func TestWalletService_GetBalance_Success(t *testing.T) {
	repo := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			return 123.45, nil
		},
	}
	ws := newWS(repo)
	bal, code := ws.GetBalance(context.Background(), "addr")
	assert.Equal(t, 123.45, bal)
	assert.Equal(t, domain.CodeOK, code)
}

func TestWalletService_GetWallet_NotFound(t *testing.T) {
	repo := &MockWalletRepository{
		GetWalletFunc: func(ctx context.Context, address string) (*domain.Wallet, error) {
			return nil, domain.ErrNotFound
		},
	}
	ws := newWS(repo)
	w, code := ws.GetWallet(context.Background(), "addr")
	assert.Nil(t, w)
	assert.Equal(t, domain.CodeWalletNotFound, code)
}

func TestWalletService_GetWallet_Internal(t *testing.T) {
	repo := &MockWalletRepository{
		GetWalletFunc: func(ctx context.Context, address string) (*domain.Wallet, error) {
			return nil, errors.New("fail")
		},
	}
	ws := newWS(repo)
	w, code := ws.GetWallet(context.Background(), "addr")
	assert.Nil(t, w)
	assert.Equal(t, domain.CodeInternal, code)
}

func TestWalletService_GetWallet_Success(t *testing.T) {
	repo := &MockWalletRepository{
		GetWalletFunc: func(ctx context.Context, address string) (*domain.Wallet, error) {
			return &domain.Wallet{Address: "addr", Balance: 100}, nil
		},
	}
	ws := newWS(repo)
	w, code := ws.GetWallet(context.Background(), "addr")
	assert.NotNil(t, w)
	assert.Equal(t, domain.CodeOK, code)
}
