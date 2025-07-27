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

func newTS(walletRepo *MockWalletRepository, txRepo *MockTransactionRepository, logger *MockLogger) *service.TransactionService {
	return service.NewTransactionService(txRepo, walletRepo, logger)
}

func TestTransactionService_GetLastTransactions_InvalidLimit(t *testing.T) {
	ts := newTS(&MockWalletRepository{}, &MockTransactionRepository{}, &MockLogger{})
	trs, code := ts.GetLastTransactions(context.Background(), 0)
	assert.Equal(t, domain.CodeInvalidLimit, code)
	assert.Nil(t, trs)
}

func TestTransactionService_GetLastTransactions_InternalError(t *testing.T) {
	tr := &MockTransactionRepository{
		GetLastTransactionsFunc: func(ctx context.Context, limit int) ([]domain.Transaction, error) {
			return nil, errors.New("fail")
		},
	}
	ts := newTS(&MockWalletRepository{}, tr, &MockLogger{})
	trs, code := ts.GetLastTransactions(context.Background(), 10)
	assert.Equal(t, domain.CodeInternal, code)
	assert.Nil(t, trs)
}

func TestTransactionService_GetLastTransactions_Success(t *testing.T) {
	tr := &MockTransactionRepository{
		GetLastTransactionsFunc: func(ctx context.Context, limit int) ([]domain.Transaction, error) {
			return []domain.Transaction{{Id: 1}}, nil
		},
	}
	ts := newTS(&MockWalletRepository{}, tr, &MockLogger{})
	trs, code := ts.GetLastTransactions(context.Background(), 10)
	assert.Equal(t, domain.CodeOK, code)
	assert.Len(t, trs, 1)
}

func TestTransactionService_GetTransactionById_NotFound(t *testing.T) {
	tr := &MockTransactionRepository{
		GetTransactionByIdFunc: func(ctx context.Context, id int64) (*domain.Transaction, error) {
			return nil, domain.ErrNotFound
		},
	}
	ts := newTS(&MockWalletRepository{}, tr, &MockLogger{})
	trn, code := ts.GetTransactionById(context.Background(), 1)
	assert.Equal(t, domain.CodeTransactionNotFound, code)
	assert.Nil(t, trn)
}

func TestTransactionService_GetTransactionById_InternalError(t *testing.T) {
	tr := &MockTransactionRepository{
		GetTransactionByIdFunc: func(ctx context.Context, id int64) (*domain.Transaction, error) {
			return nil, errors.New("fail")
		},
	}
	ts := newTS(&MockWalletRepository{}, tr, &MockLogger{})
	trn, code := ts.GetTransactionById(context.Background(), 1)
	assert.Equal(t, domain.CodeInternal, code)
	assert.Nil(t, trn)
}

func TestTransactionService_GetTransactionById_Success(t *testing.T) {
	tr := &MockTransactionRepository{
		GetTransactionByIdFunc: func(ctx context.Context, id int64) (*domain.Transaction, error) {
			return &domain.Transaction{Id: 1}, nil
		},
	}
	ts := newTS(&MockWalletRepository{}, tr, &MockLogger{})
	trn, code := ts.GetTransactionById(context.Background(), 1)
	assert.Equal(t, domain.CodeOK, code)
	assert.NotNil(t, trn)
}

func TestTransactionService_GetTransactionByInfo_NotFound(t *testing.T) {
	tr := &MockTransactionRepository{
		GetTransactionByInfoFunc: func(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, error) {
			return nil, domain.ErrNotFound
		},
	}
	ts := newTS(&MockWalletRepository{}, tr, &MockLogger{})
	trn, code := ts.GetTransactionByInfo(context.Background(), "from", "to", time.Now())
	assert.Equal(t, domain.CodeTransactionNotFound, code)
	assert.Nil(t, trn)
}

func TestTransactionService_GetTransactionByInfo_InternalError(t *testing.T) {
	tr := &MockTransactionRepository{
		GetTransactionByInfoFunc: func(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, error) {
			return nil, errors.New("fail")
		},
	}
	ts := newTS(&MockWalletRepository{}, tr, &MockLogger{})
	trn, code := ts.GetTransactionByInfo(context.Background(), "from", "to", time.Now())
	assert.Equal(t, domain.CodeInternal, code)
	assert.Nil(t, trn)
}

func TestTransactionService_GetTransactionByInfo_Success(t *testing.T) {
	tr := &MockTransactionRepository{
		GetTransactionByInfoFunc: func(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, error) {
			return &domain.Transaction{Id: 1}, nil
		},
	}
	ts := newTS(&MockWalletRepository{}, tr, &MockLogger{})
	trn, code := ts.GetTransactionByInfo(context.Background(), "from", "to", time.Now())
	assert.Equal(t, domain.CodeOK, code)
	assert.NotNil(t, trn)
} 