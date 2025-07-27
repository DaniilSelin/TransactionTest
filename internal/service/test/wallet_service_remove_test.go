package test

import (
	"context"
	"errors"
	"testing"
	"TransactionTest/internal/service"
	"TransactionTest/internal/domain"
	"github.com/stretchr/testify/assert"
)

func newWS(repo *MockWalletRepository, logger *MockLogger) *service.WalletService {
	return service.NewWalletService(repo, logger)
}

func TestWalletService_RemoveWallet_NotFound(t *testing.T) {
	repo := &MockWalletRepository{
		RemoveWalletFunc: func(ctx context.Context, address string) error {
			return domain.ErrNotFound
		},
	}
	ws := newWS(repo, &MockLogger{})
	code := ws.RemoveWallet(context.Background(), "addr")
	assert.Equal(t, domain.CodeWalletNotFound, code)
}

func TestWalletService_RemoveWallet_Internal(t *testing.T) {
	repo := &MockWalletRepository{
		RemoveWalletFunc: func(ctx context.Context, address string) error {
			return errors.New("fail")
		},
	}
	ws := newWS(repo, &MockLogger{})
	code := ws.RemoveWallet(context.Background(), "addr")
	assert.Equal(t, domain.CodeInternal, code)
}

func TestWalletService_RemoveWallet_Success(t *testing.T) {
	repo := &MockWalletRepository{
		RemoveWalletFunc: func(ctx context.Context, address string) error {
			return nil
		},
	}
	ws := newWS(repo, &MockLogger{})
	code := ws.RemoveWallet(context.Background(), "addr")
	assert.Equal(t, domain.CodeOK, code)
} 