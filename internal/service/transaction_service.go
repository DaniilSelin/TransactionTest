package service

import (
	"context"
    "errors"
    "time"

	"TransactionTest/internal/domain"
	"TransactionTest/internal/logger"
	
    "go.uber.org/zap"
)

type TransactionService struct {
    transactionRepo ITransactionRepository
    walletRepo      IWalletRepository
    log logger.Logger
}

func NewTransactionService(tr ITransactionRepository, wr IWalletRepository, l logger.Logger) *TransactionService {
    return &TransactionService{
        transactionRepo: tr,
        walletRepo:      wr,
        log: l,
    }
}

/*
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
}*/

func (ts *TransactionService) GetLastTransactions(ctx context.Context, limit int) ([]domain.Transaction, domain.ErrorCode) {
    if limit <= 0 {
		ts.log.Warn(ctx, "GetLastTransactions: limit must be greater than zero")
		return nil, domain.CodeInvalidLimit
	}

    transactions, err := ts.transactionRepo.GetLastTransactions(ctx, limit)
    if err != nil {
        ts.log.Error(ctx, "GetLastTransactions", zap.Error(err))
        return nil, domain.CodeInternal
    }
    ts.log.Info(ctx, "GetTransactionByInfo: success get last transaction", zap.Int("limit", limit))
    return transactions, domain.CodeOK
}

func (ts *TransactionService) GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, domain.ErrorCode) {
    transaction, err := ts.transactionRepo.GetTransactionById(ctx, id)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
			ts.log.Warn(ctx, "GetTransactionById: transaction not found", zap.Error(err))
			return nil, domain.CodeTransactionNotFound
		}
		ts.log.Error(ctx, "GetTransactionById",  zap.Error(err))
		return nil, domain.CodeInternal
    }
    ts.log.Info(ctx, "GetTransactionByInfo: success get transaction", zap.Int64("id", id))
    return transaction, domain.CodeOK
}

func (ts *TransactionService) GetTransactionByInfo(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, domain.ErrorCode) {
    transaction, err := ts.transactionRepo.GetTransactionByInfo(ctx, from, to, createdAt)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
			ts.log.Warn(ctx, "GetTransactionByInfo: transaction not found", zap.Error(err))
			return nil, domain.CodeTransactionNotFound
		}
		ts.log.Error(ctx, "GetTransactionByInfo",  zap.Error(err))
		return nil, domain.CodeInternal
    }
    ts.log.Info(
        ctx, 
        "GetTransactionByInfo: success get transaction",
        zap.String("from", from),
        zap.String("to", to),
        zap.Time("createdAt", createdAt),
    )
    return transaction, domain.CodeOK
}

func (ts *TransactionService) RemoveTransaction(ctx context.Context, id int64) domain.ErrorCode {
    err := ts.transactionRepo.RemoveTransaction(ctx, id)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
			ts.log.Warn(ctx, "RemoveTransaction: transaction not found", zap.Error(err))
			return domain.CodeTransactionNotFound
		}
		ts.log.Error(ctx, "RemoveTransaction",  zap.Error(err))
		return domain.CodeInternal
    }
    ts.log.Info(ctx, "RemoveTransaction: success remove transaction", zap.Int64("id", id))
    return domain.CodeOK
}
