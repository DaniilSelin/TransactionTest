package test

import (
	"context"
	"TransactionTest/internal/domain"
)

type MockTransactionService struct {
	SendMoneyFunc                func(ctx context.Context, from, to string, amount float64) domain.ErrorCode
	GetLastTransactionsFunc      func(ctx context.Context, limit int) ([]domain.Transaction, domain.ErrorCode)
	GetTransactionByIdFunc       func(ctx context.Context, id int64) (*domain.Transaction, domain.ErrorCode)
	GetTransactionByInfoFunc     func(ctx context.Context, from, to string, createdAt interface{}) (*domain.Transaction, domain.ErrorCode)
	RemoveTransactionFunc        func(ctx context.Context, id int64) domain.ErrorCode
}

func (m *MockTransactionService) SendMoney(ctx context.Context, from, to string, amount float64) domain.ErrorCode {
	return m.SendMoneyFunc(ctx, from, to, amount)
}
func (m *MockTransactionService) GetLastTransactions(ctx context.Context, limit int) ([]domain.Transaction, domain.ErrorCode) {
	return m.GetLastTransactionsFunc(ctx, limit)
}
func (m *MockTransactionService) GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, domain.ErrorCode) {
	return m.GetTransactionByIdFunc(ctx, id)
}
func (m *MockTransactionService) GetTransactionByInfo(ctx context.Context, from, to string, createdAt interface{}) (*domain.Transaction, domain.ErrorCode) {
	return m.GetTransactionByInfoFunc(ctx, from, to, createdAt)
}
func (m *MockTransactionService) RemoveTransaction(ctx context.Context, id int64) domain.ErrorCode {
	return m.RemoveTransactionFunc(ctx, id)
}

type MockWalletService struct {
	CreateWalletFunc             func(ctx context.Context, balance float64) (string, domain.ErrorCode)
	GetBalanceFunc               func(ctx context.Context, address string) (float64, domain.ErrorCode)
	GetWalletFunc                func(ctx context.Context, address string) (*domain.Wallet, domain.ErrorCode)
	UpdateBalanceFunc            func(ctx context.Context, address string, newBalance float64) domain.ErrorCode
	RemoveWalletFunc             func(ctx context.Context, address string) domain.ErrorCode
	CreateWalletsForSeedingFunc  func(ctx context.Context, count int, balance float64, failOnError bool) (<-chan string, <-chan error, bool)
}

func (m *MockWalletService) CreateWallet(ctx context.Context, balance float64) (string, domain.ErrorCode) {
	return m.CreateWalletFunc(ctx, balance)
}
func (m *MockWalletService) GetBalance(ctx context.Context, address string) (float64, domain.ErrorCode) {
	return m.GetBalanceFunc(ctx, address)
}
func (m *MockWalletService) GetWallet(ctx context.Context, address string) (*domain.Wallet, domain.ErrorCode) {
	return m.GetWalletFunc(ctx, address)
}
func (m *MockWalletService) UpdateBalance(ctx context.Context, address string, newBalance float64) domain.ErrorCode {
	return m.UpdateBalanceFunc(ctx, address, newBalance)
}
func (m *MockWalletService) RemoveWallet(ctx context.Context, address string) domain.ErrorCode {
	return m.RemoveWalletFunc(ctx, address)
}
func (m *MockWalletService) CreateWalletsForSeeding(ctx context.Context, count int, balance float64, failOnError bool) (<-chan string, <-chan error, bool) {
	return m.CreateWalletsForSeedingFunc(ctx, count, balance, failOnError)
} 