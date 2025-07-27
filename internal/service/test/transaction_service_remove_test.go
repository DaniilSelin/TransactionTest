package test

import (
	"context"
	"testing"
	"errors"
	"TransactionTest/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestTransactionService_RemoveTransaction_NotFound(t *testing.T) {
	tr := &MockTransactionRepository{
		RemoveTransactionFunc: func(ctx context.Context, id int64) error {
			return domain.ErrNotFound
		},
	}
	ts := newTS(&MockWalletRepository{}, tr)
	code := ts.RemoveTransaction(context.Background(), 1)
	assert.Equal(t, domain.CodeTransactionNotFound, code)
}

func TestTransactionService_RemoveTransaction_InternalError(t *testing.T) {
	tr := &MockTransactionRepository{
		RemoveTransactionFunc: func(ctx context.Context, id int64) error {
			return errors.New("fail")
		},
	}
	ts := newTS(&MockWalletRepository{}, tr)
	code := ts.RemoveTransaction(context.Background(), 1)
	assert.Equal(t, domain.CodeInternal, code)
}

func TestTransactionService_RemoveTransaction_Success(t *testing.T) {
	tr := &MockTransactionRepository{
		RemoveTransactionFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}
	ts := newTS(&MockWalletRepository{}, tr)
	code := ts.RemoveTransaction(context.Background(), 1)
	assert.Equal(t, domain.CodeOK, code)
}
