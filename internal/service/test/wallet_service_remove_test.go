package test

import (
	"TransactionTest/internal/domain"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalletService_RemoveWallet_NotFound(t *testing.T) {
	repo := &MockWalletRepository{
		RemoveWalletFunc: func(ctx context.Context, address string) error {
			return domain.ErrNotFound
		},
	}

	ws := newWS(repo)
	code := ws.RemoveWallet(context.Background(), "addr")
	assert.Equal(t, domain.CodeWalletNotFound, code)
}

func TestWalletService_RemoveWallet_Internal(t *testing.T) {
	repo := &MockWalletRepository{
		RemoveWalletFunc: func(ctx context.Context, address string) error {
			return errors.New("some internal error")
		},
	}

	ws := newWS(repo)
	code := ws.RemoveWallet(context.Background(), "addr")
	assert.Equal(t, domain.CodeInternal, code)
}

func TestWalletService_RemoveWallet_Success(t *testing.T) {
	repo := &MockWalletRepository{
		RemoveWalletFunc: func(ctx context.Context, address string) error {
			return nil
		},
	}

	ws := newWS(repo)
	code := ws.RemoveWallet(context.Background(), "addr")
	assert.Equal(t, domain.CodeOK, code)
}
