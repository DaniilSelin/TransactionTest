package service

import (
	"context"
	"net/http"

	"TransactionTest/internal/domain"
	"TransactionTest/internal/errors"
	"TransactionTest/internal/repository"
	goErrors "errors"
	
	_ "github.com/google/uuid"
	"time"
)

type TransactionService struct {
    transactionRepo TransactionRepositoryInterface
    walletRepo      WalletRepositoryInterface
}

func NewTransactionService(tr TransactionRepositoryInterface, wr WalletRepositoryInterface) *TransactionService {
    return &TransactionService{
        transactionRepo: tr,
        walletRepo:      wr,
    }
}
func (ts *TransactionService) SendMoney(ctx context.Context, from, to string, amount float64) error {
    if from == to {
        return errors.NewCustomError("Sender and receiver cannot be the same", http.StatusBadRequest, nil)
    }
    if amount <= 0 {
        return errors.NewCustomError("Amount must be greater than zero", http.StatusBadRequest, nil)
    }

    fromBalance, err := ts.walletRepo.GetWalletBalance(ctx, from)
    if err != nil {
        if goErrors.Is(err, repository.ErrWalletNotFound) {
            return errors.NewCustomError("Sender wallet not found", http.StatusNotFound, err)
        }
        return errors.NewCustomError("Failed to get sender balance", http.StatusInternalServerError, err)
    }
    if fromBalance < amount {
        return errors.NewCustomError("Insufficient funds", http.StatusBadRequest, nil)
    }

    toBalance, err := ts.walletRepo.GetWalletBalance(ctx, to)
    if err != nil {
        if goErrors.Is(err, repository.ErrWalletNotFound) {
            return errors.NewCustomError("Receiver wallet not found", http.StatusNotFound, err)
        }
        return errors.NewCustomError("Failed to get receiver balance", http.StatusInternalServerError, err)
    }

    newFromBalance := fromBalance - amount
    newToBalance := toBalance + amount

    // проиводим транзакцию
    err = ts.transactionRepo.ExecuteTransfer(
        ctx,
        from,
        to,
        newFromBalance,
        newToBalance,
        amount,
    )
    if err != nil {
        return errors.NewCustomError("Transaction execution failed", http.StatusInternalServerError, err)
    }
    return nil
}

func (ts *TransactionService) GetLastTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
    if limit <= 0 {
        return nil, errors.NewCustomError("Limit must be greater than zero", http.StatusBadRequest, nil)
    }

    transactions, err := ts.transactionRepo.GetLastTransactions(ctx, limit)
    if err != nil {
        return nil, errors.NewCustomError("Failed to retrieve transactions", http.StatusInternalServerError, err)
    }

    return transactions, nil
}

func (ts *TransactionService) GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, error) {
    transaction, err := ts.transactionRepo.GetTransactionById(ctx, id)
    if err != nil {
        if goErrors.Is(err, repository.ErrTransactionNotFound) {
            return nil, errors.NewCustomError("Transaction not found", http.StatusNotFound, err)
        }
        return nil, errors.NewCustomError("Failed to get transaction", http.StatusInternalServerError, err)
    }
    return transaction, nil
}

func (ts *TransactionService) GetTransactionByInfo(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, error) {
    transaction, err := ts.transactionRepo.GetTransactionByInfo(ctx, from, to, createdAt)
    if err != nil {
        if goErrors.Is(err, repository.ErrTransactionNotFound) {
            return nil, errors.NewCustomError("Transaction not found", http.StatusNotFound, err)
        }
        return nil, errors.NewCustomError("Failed to get transaction by info", http.StatusInternalServerError, err)
    }
    return transaction, nil
}

func (ts *TransactionService) RemoveTransaction(ctx context.Context, id int64) error {
    err := ts.transactionRepo.RemoveTransaction(ctx, id)
    if err != nil {
        if goErrors.Is(err, repository.ErrTransactionNotFound) {
            return errors.NewCustomError("Transaction not found", http.StatusNotFound, err)
        }
        return errors.NewCustomError("Failed to remove transaction", http.StatusInternalServerError, err)
    }
    return nil
}
