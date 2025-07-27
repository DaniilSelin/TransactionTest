package test

import (
	"context"
	"errors"
	"testing"
	"TransactionTest/internal/service"
	"TransactionTest/internal/domain"
	"github.com/stretchr/testify/assert"
)

func newTS(walletRepo *MockWalletRepository, txRepo *MockTransactionRepository, logger *MockLogger) *service.TransactionService {
	return service.NewTransactionService(txRepo, walletRepo, logger)
}

func TestTransactionService_RemoveTransaction_NotFound(t *testing.T) {
	tr := &MockTransactionRepository{
		RemoveTransactionFunc: func(ctx context.Context, id int64) error {
			return domain.ErrNotFound
		},
	}
	ts := newTS(&MockWalletRepository{}, tr, &MockLogger{})
	code := ts.RemoveTransaction(context.Background(), 1)
	assert.Equal(t, domain.CodeTransactionNotFound, code)
}

func TestTransactionService_RemoveTransaction_InternalError(t *testing.T) {
	tr := &MockTransactionRepository{
		RemoveTransactionFunc: func(ctx context.Context, id int64) error {
			return errors.New("fail")
		},
	}
	ts := newTS(&MockWalletRepository{}, tr, &MockLogger{})
	code := ts.RemoveTransaction(context.Background(), 1)
	assert.Equal(t, domain.CodeInternal, code)
}

func TestTransactionService_RemoveTransaction_Success(t *testing.T) {
	tr := &MockTransactionRepository{
		RemoveTransactionFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}
	ts := newTS(&MockWalletRepository{}, tr, &MockLogger{})
	code := ts.RemoveTransaction(context.Background(), 1)
	assert.Equal(t, domain.CodeOK, code)
} 