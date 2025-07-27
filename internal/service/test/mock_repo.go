package test

import (
	"TransactionTest/internal/domain"
	"context"
	"errors"
	"time"
)

type MockWalletRepository struct {
	BeginTXFunc               func(ctx context.Context) (domain.TxExecutor, error)
	CreateWalletFunc          func(ctx context.Context, address string, balance float64) error
	GetWalletBalanceFunc      func(ctx context.Context, address string) (float64, error)
	GetWalletFunc             func(ctx context.Context, address string) (*domain.Wallet, error)
	UpdateWalletBalanceFunc   func(ctx context.Context, address string, balance float64) error
	RemoveWalletFunc          func(ctx context.Context, address string) error
	CreateWalletTxFunc        func(ctx context.Context, tx domain.TxExecutor, address string, balance float64) error
	UpdateWalletBalanceTxFunc func(ctx context.Context, tx domain.TxExecutor, address string, balance float64) error
}

func (m *MockWalletRepository) BeginTX(ctx context.Context) (domain.TxExecutor, error) {
	return m.BeginTXFunc(ctx)
}

func (m *MockWalletRepository) CreateWallet(ctx context.Context, address string, balance float64) error {
	return m.CreateWalletFunc(ctx, address, balance)
}

func (m *MockWalletRepository) GetWalletBalance(ctx context.Context, address string) (float64, error) {
	return m.GetWalletBalanceFunc(ctx, address)
}

func (m *MockWalletRepository) GetWallet(ctx context.Context, address string) (*domain.Wallet, error) {
	return m.GetWalletFunc(ctx, address)
}

func (m *MockWalletRepository) UpdateWalletBalance(ctx context.Context, address string, balance float64) error {
	return m.UpdateWalletBalanceFunc(ctx, address, balance)
}

func (m *MockWalletRepository) RemoveWallet(ctx context.Context, address string) error {
	if m.RemoveWalletFunc != nil {
		return m.RemoveWalletFunc(ctx, address)
	}
	return nil
}

func (m *MockWalletRepository) CreateWalletTx(ctx context.Context, tx domain.TxExecutor, address string, balance float64) error {
	return m.CreateWalletTxFunc(ctx, tx, address, balance)
}

func (m *MockWalletRepository) UpdateWalletBalanceTx(ctx context.Context, tx domain.TxExecutor, address string, balance float64) error {
	return m.UpdateWalletBalanceTxFunc(ctx, tx, address, balance)
}

type MockTransactionRepository struct {
	BeginTXFunc              func(ctx context.Context) (domain.TxExecutor, error)
	CreateTransactionFunc    func(ctx context.Context, from, to string, amount float64) (int64, error)
	CreateTransactionTxFunc  func(ctx context.Context, tx domain.TxExecutor, from, to string, amount float64) (int64, error)
	GetTransactionByIdFunc   func(ctx context.Context, id int64) (*domain.Transaction, error)
	GetTransactionByInfoFunc func(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, error)
	RemoveTransactionFunc    func(ctx context.Context, id int64) error
	GetLastTransactionsFunc  func(ctx context.Context, limit int) ([]domain.Transaction, error)
}

func (m *MockTransactionRepository) BeginTX(ctx context.Context) (domain.TxExecutor, error) {
	if m.BeginTXFunc != nil {
		return m.BeginTXFunc(ctx)
	}
	return nil, errors.New("BeginTX not implemented")
}

func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, from, to string, amount float64) (int64, error) {
	return m.CreateTransactionFunc(ctx, from, to, amount)
}

func (m *MockTransactionRepository) CreateTransactionTx(ctx context.Context, tx domain.TxExecutor, from, to string, amount float64) (int64, error) {
	return m.CreateTransactionTxFunc(ctx, tx, from, to, amount)
}

func (m *MockTransactionRepository) GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, error) {
	return m.GetTransactionByIdFunc(ctx, id)
}

func (m *MockTransactionRepository) GetTransactionByInfo(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, error) {
	return m.GetTransactionByInfoFunc(ctx, from, to, createdAt)
}

func (m *MockTransactionRepository) RemoveTransaction(ctx context.Context, id int64) error {
	return m.RemoveTransactionFunc(ctx, id)
}

func (m *MockTransactionRepository) GetLastTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
	return m.GetLastTransactionsFunc(ctx, limit)
}

type MockTxExecutor struct{}

func (m *MockTxExecutor) Commit(ctx context.Context) error {
	return nil
}

func (m *MockTxExecutor) Rollback(ctx context.Context) error {
	return nil
}

func (m *MockTxExecutor) Exec(ctx context.Context, sql string, arguments ...interface{}) (domain.CommandTag, error) {
	return nil, nil
}

func (m *MockTxExecutor) QueryRow(ctx context.Context, sql string, args ...interface{}) domain.Row {
	return nil
}
