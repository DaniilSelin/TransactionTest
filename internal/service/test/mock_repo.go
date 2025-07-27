package test

import (
	"context"
	"time"
	"TransactionTest/internal/domain"
)

type MockWalletRepository struct {
	CreateWalletFunc           func(ctx context.Context, address string, balance float64) error
	GetWalletBalanceFunc       func(ctx context.Context, address string) (float64, error)
	GetWalletFunc              func(ctx context.Context, address string) (*domain.Wallet, error)
	UpdateWalletBalanceFunc    func(ctx context.Context, address string, balance float64) error
	RemoveWalletFunc           func(ctx context.Context, address string) error
	CreateWalletTxFunc         func(ctx context.Context, tx interface{}, address string, balance float64) error
	UpdateWalletBalanceTxFunc  func(ctx context.Context, tx interface{}, address string, balance float64) error
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
	return m.RemoveWalletFunc(ctx, address)
}
func (m *MockWalletRepository) CreateWalletTx(ctx context.Context, tx interface{}, address string, balance float64) error {
	return m.CreateWalletTxFunc(ctx, tx, address, balance)
}
func (m *MockWalletRepository) UpdateWalletBalanceTx(ctx context.Context, tx interface{}, address string, balance float64) error {
	return m.UpdateWalletBalanceTxFunc(ctx, tx, address, balance)
}

type MockTransactionRepository struct {
	CreateTransactionFunc      func(ctx context.Context, from, to string, amount float64) (int64, error)
	CreateTransactionTxFunc    func(ctx context.Context, tx interface{}, from, to string, amount float64) (int64, error)
	GetTransactionByIdFunc     func(ctx context.Context, id int64) (*domain.Transaction, error)
	GetTransactionByInfoFunc   func(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, error)
	RemoveTransactionFunc      func(ctx context.Context, id int64) error
	GetLastTransactionsFunc    func(ctx context.Context, limit int) ([]domain.Transaction, error)
}

func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, from, to string, amount float64) (int64, error) {
	return m.CreateTransactionFunc(ctx, from, to, amount)
}
func (m *MockTransactionRepository) CreateTransactionTx(ctx context.Context, tx interface{}, from, to string, amount float64) (int64, error) {
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