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

func (ts *TransactionService) SendMoney(ctx context.Context, from, to string, amount float64) domain.ErrorCode {
	if from == to {
		ts.log.Warn(ctx, "SendMoney: self transfer not allowed")
		return domain.CodeInvalidTransaction
	}
	if amount <= 0 {
		ts.log.Warn(ctx, "SendMoney: amount must be positive")
		return domain.CodeNegativeAmount
	}

	fromBalance, err := ts.walletRepo.GetWalletBalance(ctx, from)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ts.log.Warn(ctx, "SendMoney: sender wallet not found", zap.Error(err))
			return domain.CodeWalletNotFound
		}
		ts.log.Error(ctx, "SendMoney: failed to get sender balance", zap.Error(err))
		return domain.CodeInternal
	}
	if fromBalance < amount {
		ts.log.Warn(ctx, "SendMoney: insufficient funds", zap.Float64("balance", fromBalance), zap.Float64("amount", amount))
		return domain.CodeInsufficientFunds
	}

	toBalance, err := ts.walletRepo.GetWalletBalance(ctx, to)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ts.log.Warn(ctx, "SendMoney: receiver wallet not found", zap.Error(err))
			return domain.CodeWalletNotFound
		}
		ts.log.Error(ctx, "SendMoney: failed to get receiver balance", zap.Error(err))
		return domain.CodeInternal
	}

	newFromBalance := fromBalance - amount
	newToBalance := toBalance + amount

	tx, err := ts.transactionRepo.BeginTX(ctx)
	if err != nil {
		ts.log.Error(ctx, "SendMoney: failed to begin transaction", zap.Error(err))
		return domain.CodeInternal
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = ts.walletRepo.UpdateWalletBalanceTx(ctx, tx, from, newFromBalance)
	if err != nil {
		ts.log.Error(ctx, "SendMoney: failed to update sender balance", zap.Error(err))
		return domain.CodeInternal
	}
	err = ts.walletRepo.UpdateWalletBalanceTx(ctx, tx, to, newToBalance)
	if err != nil {
		ts.log.Error(ctx, "SendMoney: failed to update receiver balance", zap.Error(err))
		return domain.CodeInternal
	}
	_, err = ts.transactionRepo.CreateTransactionTx(ctx, tx, from, to, amount)
	if err != nil {
		ts.log.Error(ctx, "SendMoney: failed to create transaction record", zap.Error(err))
		return domain.CodeInternal
	}
	if err = tx.Commit(ctx); err != nil {
		ts.log.Error(ctx, "SendMoney: failed to commit transaction", zap.Error(err))
		return domain.CodeInternal
	}

	ts.log.Info(ctx, "SendMoney: transaction completed successfully", 
		zap.String("from", from), 
		zap.String("to", to), 
		zap.Float64("amount", amount))
	return domain.CodeOK
}

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
