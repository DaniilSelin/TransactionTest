package service

import (
	"context"
	"time"

	"TransactionTest/internal/domain"
)

type IWalletRepository interface {
	BeginTX(ctx context.Context) (domain.TxExecutor, error)
	CreateWalletTx(ctx context.Context, tx domain.TxExecutor, address string, balance float64) error
	UpdateWalletBalanceTx(ctx context.Context, tx domain.TxExecutor, address string, balance float64) error
	CreateWallet(ctx context.Context, address string, balance float64) error
	GetWalletBalance(ctx context.Context, address string) (float64, error)
	GetWallet(ctx context.Context, address string) (*domain.Wallet, error)
	UpdateWalletBalance(ctx context.Context, address string, balance float64) error
	RemoveWallet(ctx context.Context, address string) error
}

type ITransactionRepository interface {
	BeginTX(ctx context.Context) (domain.TxExecutor, error)
	CreateTransactionTx(ctx context.Context, tx domain.TxExecutor, from, to string, amount float64) (int64, error)
	CreateTransaction(ctx context.Context, from, to string, amount float64) (int64, error)
	GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, error)
	GetTransactionByInfo(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, error)
	RemoveTransaction(ctx context.Context, id int64) error
	GetLastTransactions(ctx context.Context, limit int) ([]domain.Transaction, error)
}
