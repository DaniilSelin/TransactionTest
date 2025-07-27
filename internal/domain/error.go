package domain

import (
	"errors"
)

// ошибки repository
var (
    ErrNotFound        = errors.New("not found")
    ErrConflict        = errors.New("conflict")
    ErrInvalidInput    = errors.New("invalid input")
    ErrInternal        = errors.New("internal server error")
)

//  Ошибки кошельков 
var (
    ErrNegativeBalance     = errors.New("negative balance not allowed")
    ErrWalletAlreadyExists = errors.New("wallet address already exists")
    ErrInsufficientFunds   = errors.New("insufficient funds")
)

// Ошибки транзакций
var (
    ErrNegativeAmount      = errors.New("amount must be positive")
    ErrInvalidTransaction  = errors.New("invalid transaction")
    ErrTransactionConflict = errors.New("transaction conflict")
    ErrSelfTransfer        = errors.New("cannot transfer to self")
)

// коды ошибок слоя бизнесс логики
type ErrorCode string

const (
    CodeOK                 ErrorCode = ""   
    CodeWalletNotFound   ErrorCode = "WALLET_NOT_FOUND"
    CodeInsufficientFunds ErrorCode = "INSUFFICIENT_FUNDS"
    CodeDuplicateWallet   ErrorCode = "DUPLICATE_WALLET"
    CodeNegativeBalance    ErrorCode = "NEGATIVE_BALANCE"
    CodeInternal           ErrorCode = "INTERNAL_ERROR"
    CodeInvalidLimit          ErrorCode = "INVALID_LIMIT"
    CodeTransactionNotFound   ErrorCode = "TRANSACTION_NOT_FOUND"
    CodeNegativeAmount        ErrorCode = "NEGATIVE_AMOUNT"
    CodeInvalidTransaction    ErrorCode = "INVALID_TRANSACTION"
    CodeInvalidRequestBody      ErrorCode = "INVALID_REQUEST_BODY"
)
