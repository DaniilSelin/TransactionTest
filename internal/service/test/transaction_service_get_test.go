package test

import (
	"context"
	"errors"
	"testing"
	"time"
	"TransactionTest/internal/domain"
	"TransactionTest/internal/logger"
	"TransactionTest/internal/service"
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func newTS(walletRepo service.IWalletRepository, txRepo service.ITransactionRepository) *service.TransactionService {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	log, _ := logger.New(&cfg)
    return service.NewTransactionService(txRepo, walletRepo, log)
}

func newWS(walletRepo service.IWalletRepository) *service.WalletService {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	log, _ := logger.New(&cfg)
	return service.NewWalletService(walletRepo, log)
}

func TestTransactionService_GetLastTransactions_InvalidLimit(t *testing.T) {
	ts := newTS(&MockWalletRepository{}, &MockTransactionRepository{})
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
	ts := newTS(&MockWalletRepository{}, tr)
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
	ts := newTS(&MockWalletRepository{}, tr)
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
	ts := newTS(&MockWalletRepository{}, tr)
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
	ts := newTS(&MockWalletRepository{}, tr)
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
	ts := newTS(&MockWalletRepository{}, tr)
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
	ts := newTS(&MockWalletRepository{}, tr)
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
	ts := newTS(&MockWalletRepository{}, tr)
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
	ts := newTS(&MockWalletRepository{}, tr)
	trn, code := ts.GetTransactionByInfo(context.Background(), "from", "to", time.Now())
	assert.Equal(t, domain.CodeOK, code)
	assert.NotNil(t, trn)
} 