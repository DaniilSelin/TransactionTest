package test

import (
	"context"
	"errors"
	"testing"
	"time"
	"TransactionTest/internal/service"
	"TransactionTest/internal/domain"
	"github.com/stretchr/testify/assert"
)

func newWS(repo *MockWalletRepository, logger *MockLogger) *service.WalletService {
	return service.NewWalletService(repo, logger)
}

func TestWalletService_CreateWalletsForSeeding_BeginError(t *testing.T) {
	repo := &MockWalletRepository{
		BeginFunc: func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("fail-begin")
		},
	}
	ws := newWS(repo, &MockLogger{})
	done, errChan, ok := ws.CreateWalletsForSeeding(context.Background(), 2, 100, false)
	assert.False(t, ok)
	assert.NotNil(t, done)
	assert.NotNil(t, errChan)
}

func TestWalletService_CreateWalletsForSeeding_Success(t *testing.T) {
	calls := 0
	repo := &MockWalletRepository{
		BeginFunc: func(ctx context.Context) (interface{}, error) {
			return &struct{}{}, nil
		},
		CreateWalletTxFunc: func(ctx context.Context, tx interface{}, address string, balance float64) error {
			calls++
			return nil
		},
	}
	ws := newWS(repo, &MockLogger{})
	done, errChan, ok := ws.CreateWalletsForSeeding(context.Background(), 2, 100, false)
	assert.True(t, ok)
	for range done {
	}
	for err := range errChan {
		assert.NoError(t, err)
	}
	assert.Equal(t, 2, calls)
}

func TestWalletService_CreateWalletsForSeeding_CreateError(t *testing.T) {
	repo := &MockWalletRepository{
		BeginFunc: func(ctx context.Context) (interface{}, error) {
			return &struct{}{}, nil
		},
		CreateWalletTxFunc: func(ctx context.Context, tx interface{}, address string, balance float64) error {
			return errors.New("fail-create")
		},
	}
	ws := newWS(repo, &MockLogger{})
	done, errChan, ok := ws.CreateWalletsForSeeding(context.Background(), 2, 100, false)
	assert.True(t, ok)
	for range done {
	}
	errCount := 0
	for err := range errChan {
		assert.Error(t, err)
		errCount++
	}
	assert.Equal(t, 2, errCount)
}

func TestWalletService_CreateWalletsForSeeding_FailOnError(t *testing.T) {
	repo := &MockWalletRepository{
		BeginFunc: func(ctx context.Context) (interface{}, error) {
			return &struct{}{}, nil
		},
		CreateWalletTxFunc: func(ctx context.Context, tx interface{}, address string, balance float64) error {
			return errors.New("fail-create")
		},
	}
	ws := newWS(repo, &MockLogger{})
	done, errChan, ok := ws.CreateWalletsForSeeding(context.Background(), 2, 100, true)
	assert.True(t, ok)
	for range done {
	}
	errCount := 0
	for err := range errChan {
		assert.Error(t, err)
		errCount++
	}
	assert.Equal(t, 1, errCount)
} 