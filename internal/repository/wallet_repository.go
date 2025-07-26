package repository

import (
    "context"
    "fmt"

    "TransactionTest/internal/domain"
)

type WalletRepository struct {
	db IDB
}

func NewWalletRepository(db IDB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (wr *WalletRepository) CreateWallet(ctx context.Context, address string, balance float64) error {
	query := `INSERT INTO wallets (address, balance) VALUES ($1, $2)`

	_, err := wr.db.Exec(ctx, query, address, balance)
	if err != nil {
		if dbErr, ok := err.(DBError); ok {
			if dbErr.SQLState() == ErrCodeCheckViolation && dbErr.ConstraintName() == ConstraintBalanceNonNegative {
				return domain.ErrNegativeBalance
			}
			if dbErr.SQLState() == ErrCodeUniqueViolation {
				return domain.ErrWalletAlreadyExists
			}
		}
		return fmt.Errorf("%w: failed to create wallet: %w", domain.ErrInternal, err)
	}

	return nil
}

func (wr *WalletRepository) GetWalletBalance(ctx context.Context, address string) (float64, error) {
	query := `SELECT balance FROM wallets WHERE address = $1`

	var balance float64

	err := wr.db.QueryRow(ctx, query, address).Scan(&balance)
	if err != nil {
        if dbErr, ok := err.(DBError); ok && dbErr.SQLState() == "no_rows" {
        	return 0, domain.ErrNotFound
        }
        return 0, fmt.Errorf("%w: failed to find wallet %v: %w",  domain.ErrInternal, address, err)
	}

	return balance, nil
}

func (wr *WalletRepository) GetWallet(ctx context.Context, address string) (*domain.Wallet, error) {
    query := `SELECT address, balance, created_at 
    		  FROM wallets WHERE address = $1`

    var w domain.Wallet

    err := wr.db.QueryRow(ctx, query, address).Scan(
    	&w.Address,
    	&w.Balance, 
		&w.CreatedAt,
    )

    if err != nil {
        if dbErr, ok := err.(DBError); ok && dbErr.SQLState() == "no_rows" {
        	return nil, domain.ErrNotFound
        }
        return nil, fmt.Errorf("%w: failed to find wallet %v: %w", domain.ErrInternal, address, err)
    }

    return &w, nil
}

func (wr *WalletRepository) UpdateWalletBalance(ctx context.Context, address string, balance float64) error {
	query := `UPDATE wallets SET balance = $1 WHERE address = $2`

	result, err := wr.db.Exec(ctx, query, balance, address)
    if err != nil {
        if dbErr, ok := err.(DBError); ok {
        	if dbErr.SQLState() == ErrCodeCheckViolation && dbErr.ConstraintName() == ConstraintBalanceNonNegative {
        		return domain.ErrNegativeBalance
        	}
        }
        return fmt.Errorf("%w: failed to update wallet %v: %w", domain.ErrInternal, address, err)
    }

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (wr *WalletRepository) RemoveWallet(ctx context.Context, address string) error {
	query := `DELETE FROM wallets WHERE address = $1`

	result, err := wr.db.Exec(ctx, query, address)
	if err != nil {
		return fmt.Errorf("%w: failed to delete wallet %v: %w", domain.ErrInternal, address, err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (wr *WalletRepository) BatchCreateWallets(
	ctx context.Context,
	failOnError bool,
	wallets  <-chan domain.Wallet,
	done chan<- string,
	errChan chan<- error,
) {
	tx, err := wr.db.Begin(ctx)
	if err != nil {
	    errChan <- err
	    close(done)
	    close(errChan)
	    return
	}

	query := `INSERT INTO wallets (address, balance) VALUES ($1, $2)`

	for w := range wallets {
		select {
		case <-ctx.Done():
		    tx.Rollback(ctx)
		    close(done)
		    close(errChan)
		    return
		default:
		}
	
		if _, err := tx.Exec(ctx, query, w.Address, w.Balance); err != nil {
		    if failOnError {
				tx.Rollback(ctx)
				errChan <- fmt.Errorf("insert %s: %w", w.Address, err)
				close(done)
				close(errChan)
				return
		    }
		    errChan <- fmt.Errorf("WARN: insert failed for %s: %v", w.Address, err)
		    continue
		}
		done <- w.Address
	}
	
	if err := tx.Commit(ctx); err != nil {
		errChan <- fmt.Errorf("commit: %w", err)
	}
	close(done)
	close(errChan)
}