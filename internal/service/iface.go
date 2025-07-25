package service

import (
	"context"
	"time"

	"TransactionTest/internal/domain"
)

type IWalletRepository interface {
	CreateWallet(ctx context.Context, address string, balance float64) error
	GetWalletBalance(ctx context.Context, address string) (float64, error)
	GetWallet(ctx context.Context, address string) (*domain.Wallet, error)
	UpdateWalletBalabnce(ctx context.Context, address string, balance float64) error
	RemoveWallet(ctx context.Context, address string) error
	IsEmpty(ctx context.Context) (bool, error)
	BatchCreateWallets(
		ctx context.Context,
		failOnError bool,
		wallets <-chan domain.Wallet,
		done chan<- string,
		errChan chan<- error,
	)
}

type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, from, to string, amount float64) (int64, error)
	ExecuteTransfer(ctx context.Context, from, to string, balance_from, balance_to, amount float64) error
	GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, error)
	GetTransactionByInfo(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, error)
	RemoveTransaction(ctx context.Context, id int64) error
	GetLastTransactions(ctx context.Context, limit int) ([]domain.Transaction, error)
}
