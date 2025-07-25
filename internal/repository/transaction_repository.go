package repository

import (
    "context"
    "fmt"
    "time"

    "TransactionTest/internal/domain"
)

type TransactionRepository struct {
	db IDB
}

func NewTransactionRepository(db IDB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (tr *TransactionRepository) CreateTransaction(ctx context.Context, from, to string, amount float64) (int64, error) {
	query := `INSERT INTO transactions (from_wallet, to_wallet, amount) 
              VALUES ($1, $2, $3) RETURNING id`

	var transactionId int64

	err := tr.db.QueryRow(ctx, query, from, to, amount).Scan(&transactionId)
	if err != nil {
		if dbErr, ok := err.(DBError); ok {
			if dbErr.SQLState() == ErrCodeCheckViolation && dbErr.ConstraintName() == ConstraintAmountPositive {
				return 0, domain.ErrNegativeAmount
			}
			if dbErr.SQLState() == ErrCodeCheckViolation && dbErr.ConstraintName() == ConstraintNoSelfTransfer {
				return 0, domain.ErrSelfTransfer
			}
			if dbErr.SQLState() == ErrCodeForeignKeyViolation {
				return 0, fmt.Errorf("%w: wallet %s or %s", domain.ErrNotFound, from, to)
			}
		}
		return 0, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return transactionId, nil
}

func (tr *TransactionRepository) ExecuteTransfer(ctx context.Context, from, to string, balance_from, balance_to, amount float64) error {
	tx, err := tr.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: transaction start failed: %w", domain.ErrInternal, err)
	}
	defer tx.Rollback(ctx)

    updateQuery := `
        UPDATE wallets 
        SET balance = CASE 
            WHEN address = $1 THEN CAST($3 AS numeric)
            WHEN address = $2 THEN CAST($4 AS numeric) 
        END 
        WHERE address IN ($1, $2)`

	_, err = tx.Exec(ctx, updateQuery, from, to, balance_from, balance_to)
    if err != nil {
        if dbErr, ok := err.(DBError); ok {
        	if dbErr.SQLState() == ErrCodeCheckViolation && dbErr.ConstraintName() == ConstraintBalanceNonNegative {
        		return domain.ErrNegativeBalance
        	}
        }
        return fmt.Errorf("balance update failed: %w", err)
    }

	_, err = tx.Exec(ctx,
		`INSERT INTO transactions 
		(from_wallet, to_wallet, amount) 
		VALUES ($1, $2, $3)`,
		from, to, amount,
	)
	if err != nil {
		if dbErr, ok := err.(DBError); ok {
			if dbErr.SQLState() == ErrCodeCheckViolation && dbErr.ConstraintName() == ConstraintAmountPositive {
				return fmt.Errorf("amount must be positive: %w", err)
			}
			if dbErr.SQLState() == ErrCodeCheckViolation && dbErr.ConstraintName() == ConstraintNoSelfTransfer {
				return fmt.Errorf("cannot transfer to self: %w", err)
			}
			if dbErr.SQLState() == ErrCodeForeignKeyViolation {
				return fmt.Errorf("wallet not found: %w", err)
			}
		}
		return fmt.Errorf("transaction record failed: %w", err)
	}

	return tx.Commit(ctx)
}

func (tr *TransactionRepository) GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, error) {
    query := `SELECT id, from_wallet, to_wallet, amount, created_at 
    		  FROM transactions WHERE id = $1`

    var t domain.Transaction

    err := tr.db.QueryRow(ctx, query, id).Scan(
        &t.Id,
        &t.From,
        &t.To,
        &t.Amount,
        &t.CreatedAt,
    )

    if err != nil {
        if dbErr, ok := err.(DBError); ok && dbErr.SQLState() == "no_rows" {
        	return nil, domain.ErrNotFound
        }
        return nil, fmt.Errorf("%w: failed to find transaction with id %v: %w", domain.ErrInternal, id, err)
    }

    return &t, nil
}
	
func (tr *TransactionRepository) GetTransactionByInfo(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, error) {
    query := `SELECT id, from_wallet, to_wallet, amount, created_at
              FROM transactions 
              WHERE from_wallet = $1 AND to_wallet = $2 AND created_at = $3`

    var t domain.Transaction

    err := tr.db.QueryRow(ctx, query, from, to, createdAt).Scan(
        &t.Id,
        &t.From,
        &t.To,
        &t.Amount,
        &t.CreatedAt,
    )

    if err != nil {
        if dbErr, ok := err.(DBError); ok && dbErr.SQLState() == "no_rows" {
        	return nil, domain.ErrNotFound
        }
        return nil, fmt.Errorf("%w: fail than look for transaction for from_wallet %v, to_wallet %v at %v: %w", domain.ErrInternal, from, to, createdAt, err)
    }

    return &t, nil
}

func (tr *TransactionRepository) RemoveTransaction(ctx context.Context, id int64) error {
	query := `DELETE FROM transactions WHERE id = $1`

	result, err := tr.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: failed to delete transaction with id %v: %w", domain.ErrInternal, id, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (tr *TransactionRepository) GetLastTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
    query := `SELECT id, from_wallet, to_wallet, amount, created_at 
              FROM transactions 
              ORDER BY created_at DESC
              LIMIT $1`

    rows, err := tr.db.Query(ctx, query, limit)
    if err != nil {
        return nil, fmt.Errorf("%w: failed to get last transactions: %w", domain.ErrInternal, err)
    }
    defer rows.Close()

    transactions := make([]domain.Transaction, 0, limit)

    for rows.Next() {
        var t domain.Transaction

        if err := rows.Scan(&t.Id, &t.From, &t.To, &t.Amount, &t.CreatedAt); err != nil {
            return nil, fmt.Errorf("%w: failed to scan transaction: %w", domain.ErrInternal, err)
        }

        transactions = append(transactions, t)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%w: error while fetching rows: %w", domain.ErrInternal, err)
    }

    return transactions, nil
}
