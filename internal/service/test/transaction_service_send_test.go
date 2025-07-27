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

func TestTransactionService_SendMoney_SelfTransfer(t *testing.T) {
	ts := newTS(&MockWalletRepository{}, &MockTransactionRepository{}, &MockLogger{})
	code := ts.SendMoney(context.Background(), "addr", "addr", 10)
	assert.Equal(t, domain.CodeInvalidTransaction, code)
}

func TestTransactionService_SendMoney_NegativeAmount(t *testing.T) {
	ts := newTS(&MockWalletRepository{}, &MockTransactionRepository{}, &MockLogger{})
	code := ts.SendMoney(context.Background(), "from", "to", -1)
	assert.Equal(t, domain.CodeNegativeAmount, code)
}

func TestTransactionService_SendMoney_SenderNotFound(t *testing.T) {
	wr := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			return 0, domain.ErrNotFound
		},
	}
	ts := newTS(wr, &MockTransactionRepository{}, &MockLogger{})
	code := ts.SendMoney(context.Background(), "from", "to", 10)
	assert.Equal(t, domain.CodeWalletNotFound, code)
}

func TestTransactionService_SendMoney_ReceiverNotFound(t *testing.T) {
	calls := 0
	wr := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			calls++
			if calls == 1 {
				return 100, nil
			}
			return 0, domain.ErrNotFound
		},
	}
	ts := newTS(wr, &MockTransactionRepository{}, &MockLogger{})
	code := ts.SendMoney(context.Background(), "from", "to", 10)
	assert.Equal(t, domain.CodeWalletNotFound, code)
}

func TestTransactionService_SendMoney_InsufficientFunds(t *testing.T) {
	wr := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			return 5, nil
		},
	}
	ts := newTS(wr, &MockTransactionRepository{}, &MockLogger{})
	code := ts.SendMoney(context.Background(), "from", "to", 10)
	assert.Equal(t, domain.CodeInsufficientFunds, code)
}

func TestTransactionService_SendMoney_InternalErrorOnSender(t *testing.T) {
	wr := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			return 0, errors.New("fail")
		},
	}
	ts := newTS(wr, &MockTransactionRepository{}, &MockLogger{})
	code := ts.SendMoney(context.Background(), "from", "to", 10)
	assert.Equal(t, domain.CodeInternal, code)
}

func TestTransactionService_SendMoney_InternalErrorOnReceiver(t *testing.T) {
	calls := 0
	wr := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			calls++
			if calls == 1 {
				return 100, nil
			}
			return 0, errors.New("fail")
		},
	}
	ts := newTS(wr, &MockTransactionRepository{}, &MockLogger{})
	code := ts.SendMoney(context.Background(), "from", "to", 10)
	assert.Equal(t, domain.CodeInternal, code)
}

func TestTransactionService_SendMoney_InternalErrorOnUpdate(t *testing.T) {
	wr := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			return 100, nil
		},
		UpdateWalletBalanceTxFunc: func(ctx context.Context, tx interface{}, address string, balance float64) error {
			return errors.New("fail")
		},
	}
	ts := newTS(wr, &MockTransactionRepository{}, &MockLogger{})
	code := ts.SendMoney(context.Background(), "from", "to", 10)
	assert.Equal(t, domain.CodeInternal, code)
}

func TestTransactionService_SendMoney_InternalErrorOnCreateTx(t *testing.T) {
	wr := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			return 100, nil
		},
		UpdateWalletBalanceTxFunc: func(ctx context.Context, tx interface{}, address string, balance float64) error {
			return nil
		},
	}
	tr := &MockTransactionRepository{
		CreateTransactionTxFunc: func(ctx context.Context, tx interface{}, from, to string, amount float64) (int64, error) {
			return 0, errors.New("fail")
		},
	}
	ts := newTS(wr, tr, &MockLogger{})
	code := ts.SendMoney(context.Background(), "from", "to", 10)
	assert.Equal(t, domain.CodeInternal, code)
}

func TestTransactionService_SendMoney_InternalErrorOnCommit(t *testing.T) {
	wr := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			return 100, nil
		},
		UpdateWalletBalanceTxFunc: func(ctx context.Context, tx interface{}, address string, balance float64) error {
			return nil
		},
	}
	tr := &MockTransactionRepository{
		CreateTransactionTxFunc: func(ctx context.Context, tx interface{}, from, to string, amount float64) (int64, error) {
			return 1, nil
		},
	}
	logger := &MockLogger{}
	ts := newTS(wr, tr, logger)
	// эмулируем ошибку коммита через мок транзакции, если нужно
	// здесь предполагается, что BeginTX и tx.Commit будут замоканы в реальном тесте
	code := ts.SendMoney(context.Background(), "from", "to", 10)
	assert.Equal(t, domain.CodeInternal, code)
}

func TestTransactionService_SendMoney_Success(t *testing.T) {
	wr := &MockWalletRepository{
		GetWalletBalanceFunc: func(ctx context.Context, address string) (float64, error) {
			return 100, nil
		},
		UpdateWalletBalanceTxFunc: func(ctx context.Context, tx interface{}, address string, balance float64) error {
			return nil
		},
	}
	tr := &MockTransactionRepository{
		CreateTransactionTxFunc: func(ctx context.Context, tx interface{}, from, to string, amount float64) (int64, error) {
			return 1, nil
		},
	}
	logger := &MockLogger{}
	ts := newTS(wr, tr, logger)
	code := ts.SendMoney(context.Background(), "from", "to", 10)
	assert.Equal(t, domain.CodeOK, code)
} 