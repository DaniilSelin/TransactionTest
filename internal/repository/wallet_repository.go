package repository

import (
    "context"
    "fmt"
    "errors"

    "TransactionSystem/internal/domain"

    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/jackc/pgx/v4"
)

var ErrWalletNotFound = errors.New("wallet not found")

// Менеджер для кошельков
type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{db: db}
}

func (wr *WalletRepository) CreateWallet(ctx context.Context, address string, balance float64) error {
	query := `INSERT INTO "TransactionSystem".wallets (address, balance) VALUES ($1, $2)`

	_, err := wr.db.Exec(ctx, query, address, balance)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

func (wr *WalletRepository) GetWalletBalance(ctx context.Context, address string) (float64, error) {
	query := `SELECT balance FROM "TransactionSystem".wallets WHERE address = $1`

	var balance float64
	
	err := wr.db.QueryRow(ctx, query, address).Scan(&balance)
	if err != nil {
		if err == pgx.ErrNoRows {
            return 0, ErrWalletNotFound
        }
        return 0, fmt.Errorf("failed to find wallet with address %v: %w", address, err)
	}

	return balance, nil
}

func (wr *WalletRepository) GetWallet(ctx context.Context, address string) (*domain.Wallet, error) {
    query := `SELECT address, balance, created_at 
    		  FROM "TransactionSystem".wallets WHERE address = $1`

    var w domain.Wallet

    err := wr.db.QueryRow(ctx, query, address).Scan(
    	&w.Address,
	    &w.Balance, 
		&w.CreatedAt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, ErrWalletNotFound
        }
        return nil, fmt.Errorf("failed to find wallet with address %v: %w", address, err)
    }

    return &w, nil
}

func (wr *WalletRepository) UpdateWalletBalabnce(ctx context.Context, address string, balance float64) error {
	query := `UPDATE "TransactionSystem".wallets SET balance = $1 WHERE address = $2`

	result, err := wr.db.Exec(ctx, query, balance, address)
    if err != nil {
        return fmt.Errorf("failed to update wallet with address %v: %w", address, err)
    }

	// Проверяем, были ли затронуты строки
	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrWalletNotFound
	}

	return nil
}

func (wr *WalletRepository) RemoveWallet(ctx context.Context, address string) error {
	query := `DELETE FROM "TransactionSystem".wallets WHERE address = $1`

	result, err := wr.db.Exec(ctx, query, address)
	if err != nil {
		return fmt.Errorf("failed to delete wallet with address %v: %w", address, err)
	}

	//  проверяем, были ли затронуты строки
	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrWalletNotFound
	}

	return nil
}

func (wr *WalletRepository) IsEmpty(ctx context.Context) (bool, error) {
    query := `SELECT COUNT(*) FROM "TransactionSystem".wallets`
    
    var count int
    err := wr.db.QueryRow(ctx, query).Scan(&count)
    if err != nil {
        return false, fmt.Errorf("failed to check wallets: %w", err)
    }
    
    return count == 0, nil
}