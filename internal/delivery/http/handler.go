package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"TransactionTest/internal/domain"
	"TransactionTest/internal/logger"
	"TransactionTest/internal/service"

	"go.uber.org/zap"
)

type ITransactionService interface {
	SendMoney(ctx context.Context, from, to string, amount float64) domain.ErrorCode
	GetLastTransactions(ctx context.Context, limit int) ([]domain.Transaction, domain.ErrorCode)
	GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, domain.ErrorCode)
	GetTransactionByInfo(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, domain.ErrorCode)
	RemoveTransaction(ctx context.Context, id int64) domain.ErrorCode
}

type IWalletService interface {
	CreateWallet(ctx context.Context, balance float64) (string, domain.ErrorCode)
	GetBalance(ctx context.Context, address string) (float64, domain.ErrorCode)
	GetWallet(ctx context.Context, address string) (*domain.Wallet, domain.ErrorCode)
	UpdateBalance(ctx context.Context, address string, newBalance float64) domain.ErrorCode
	RemoveWallet(ctx context.Context, address string) domain.ErrorCode
}

type Handler struct {
	transactionService ITransactionService
	walletService      IWalletService
	log                logger.Logger
}

func NewHandler(ts ITransactionService, ws IWalletService, l logger.Logger) *Handler {
	return &Handler{
		transactionService: ts,
		walletService:      ws,
		log:                l,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

func (h *Handler) writeJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error(ctx, "Failed to encode JSON response", zap.Error(err))
	}
}

func (h *Handler) writeError(w http.ResponseWriter, statusCode int, code domain.ErrorCode, message string) {
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Code:    string(code),
		Message: message,
	}
	h.writeJSON(w, statusCode, response)
}

func (h *Handler) handleServiceError(ctx context.Context, w http.ResponseWriter, code domain.ErrorCode, operation string) {
	switch code {
	case domain.CodeOK:
		return
	case domain.CodeWalletNotFound:
		h.log.Warn(ctx, operation+": wallet not found")
		h.writeError(w, http.StatusNotFound, code, "Wallet not found")
	case domain.CodeTransactionNotFound:
		h.log.Warn(ctx, operation+": transaction not found")
		h.writeError(w, http.StatusNotFound, code, "Transaction not found")
	case domain.CodeInsufficientFunds:
		h.log.Warn(ctx, operation+": insufficient funds")
		h.writeError(w, http.StatusBadRequest, code, "Insufficient funds")
	case domain.CodeDuplicateWallet:
		h.log.Warn(ctx, operation+": duplicate wallet")
		h.writeError(w, http.StatusConflict, code, "Wallet already exists")
	case domain.CodeNegativeBalance:
		h.log.Warn(ctx, operation+": negative balance")
		h.writeError(w, http.StatusBadRequest, code, "Negative balance not allowed")
	case domain.CodeNegativeAmount:
		h.log.Warn(ctx, operation+": negative amount")
		h.writeError(w, http.StatusBadRequest, code, "Amount must be positive")
	case domain.CodeInvalidTransaction:
		h.log.Warn(ctx, operation+": invalid transaction")
		h.writeError(w, http.StatusBadRequest, code, "Invalid transaction")
	case domain.CodeInvalidLimit:
		h.log.Warn(ctx, operation+": invalid limit")
		h.writeError(w, http.StatusBadRequest, code, "Invalid limit parameter")
	case domain.CodeInternal:
		h.log.Error(ctx, operation+": internal error")
		h.writeError(w, http.StatusInternalServerError, code, "Internal server error")
	default:
		h.log.Error(ctx, operation+": unknown error code", zap.String("code", string(code)))
		h.writeError(w, http.StatusInternalServerError, domain.CodeInternal, "Internal server error")
	}
} 