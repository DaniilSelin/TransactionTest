package seeder

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"TransactionTest/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateFunc struct {
	mock.Mock
}

func (m *mockCreateFunc) create(
	ctx context.Context,
	count int,
	balance float64,
	failOnError bool,
) (<-chan string, <-chan error, bool) {
	args := m.Called(ctx, count, balance, failOnError)
	
	doneChan := args.Get(0).(chan string)
	errChan := args.Get(1).(chan error)
	
	return (<-chan string)(doneChan), (<-chan error)(errChan), true
}

func TestSeedWallets_Success(t *testing.T) {
	mockCreator := new(mockCreateFunc)
	markerFile := filepath.Join(os.TempDir(), fmt.Sprintf("marker_%d", time.Now().UnixNano()))
	defer os.Remove(markerFile)

	cfg := config.WalletsSeedConfig{
		Enabled:    true,
		Count:      3,
		Balance:    100.0,
		MarkerFile: markerFile,
	}

	doneChan := make(chan string, cfg.Count)
	errChan := make(chan error, 1)

	for i := 0; i < cfg.Count; i++ {
		doneChan <- fmt.Sprintf("wallet_%d", i)
	}
	close(doneChan)

	mockCreator.On(
		"create",
		mock.AnythingOfType("*context.cancelCtx"),
		cfg.Count,
		cfg.Balance,
		cfg.FailOnError,
	).Return(doneChan, errChan)

	err := SeedWallets(context.Background(), cfg, mockCreator.create, func(_ context.Context, _ error) {})
	assert.NoError(t, err)
	mockCreator.AssertExpectations(t)

	content, err := os.ReadFile(markerFile)
	assert.NoError(t, err)
	assert.Equal(t, "wallet_0\nwallet_1\nwallet_2\n", string(content))
}

func TestSeedWallets_ErrorFromCreator(t *testing.T) {
	mockCreator := new(mockCreateFunc)
	markerFile := filepath.Join(os.TempDir(), fmt.Sprintf("test_marker_%d", time.Now().UnixNano()))
	defer os.Remove(markerFile)

	cfg := config.WalletsSeedConfig{
		Enabled:    true,
		Count:      1,
		Balance:    100.0,
		MarkerFile: markerFile,
	}

	doneChan := make(chan string)
	errChan := make(chan error, 1)
	errChan <- errors.New("creation error")


	mockCreator.On(
		"create",
		mock.AnythingOfType("*context.cancelCtx"),
		cfg.Count,
		cfg.Balance,
		cfg.FailOnError,
	).Return(doneChan, errChan).Once()

	err := SeedWallets(context.Background(), cfg, mockCreator.create, func(_ context.Context, _ error) {})

	close(doneChan)
	close(errChan)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "creation error")
	mockCreator.AssertExpectations(t)
}

func TestSeedWallets_WriteError(t *testing.T) {
	mockCreator := new(mockCreateFunc)
	invalidPath := filepath.Join("/invalid", fmt.Sprintf("test_%d", time.Now().UnixNano()))

	cfg := config.WalletsSeedConfig{
		Enabled:    true,
		Count:      1,
		Balance:    100.0,
		MarkerFile: invalidPath,
	}

	doneChan := make(chan string, 1)
	doneChan <- "test_wallet"
	errChan := make(chan error, 1)

	mockCreator.On(
		"create",
		mock.AnythingOfType("*context.cancelCtx"),
		cfg.Count,
		cfg.Balance,
		cfg.FailOnError,
	).Return(doneChan, errChan).Once()

	close(doneChan)
	close(errChan)
	
	err := SeedWallets(context.Background(), cfg, mockCreator.create, func(_ context.Context, _ error) {})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marker")
}