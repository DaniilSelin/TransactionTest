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

func TestWalletService_UpdateBalance_Negative(t *testing.T) {
	ws := newWS(&MockWalletRepository{}, &MockLogger{})
	code := ws.UpdateBalance(context.Background(), "addr", -1)
	assert.Equal(t, domain.CodeNegativeBalance, code)
}

func TestWalletService_UpdateBalance_NotFound(t *testing.T) {
	repo := &MockWalletRepository{
		UpdateWalletBalanceFunc: func(ctx context.Context, address string, balance float64) error {
			return domain.ErrNotFound
		},
	}
	ws := newWS(repo, &MockLogger{})
	code := ws.UpdateBalance(context.Background(), "addr", 100)
	assert.Equal(t, domain.CodeWalletNotFound, code)
}

func TestWalletService_UpdateBalance_Internal(t *testing.T) {
	repo := &MockWalletRepository{
		UpdateWalletBalanceFunc: func(ctx context.Context, address string, balance float64) error {
			return errors.New("fail")
		},
	}
	ws := newWS(repo, &MockLogger{})
	code := ws.UpdateBalance(context.Background(), "addr", 100)
	assert.Equal(t, domain.CodeInternal, code)
}

func TestWalletService_UpdateBalance_Success(t *testing.T) {
	repo := &MockWalletRepository{
		UpdateWalletBalanceFunc: func(ctx context.Context, address string, balance float64) error {
			return nil
		},
	}
	ws := newWS(repo, &MockLogger{})
	code := ws.UpdateBalance(context.Background(), "addr", 100)
	assert.Equal(t, domain.CodeOK, code)
} 