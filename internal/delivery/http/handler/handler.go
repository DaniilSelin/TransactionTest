// Package handler предоставляет HTTP обработчики для API транзакций и кошельков.
//
// Пакет содержит все HTTP обработчики для работы с:
// - Транзакциями (отправка денег, получение истории, удаление)
// - Кошельками (создание, получение баланса, обновление, удаление)
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"TransactionTest/internal/domain"
	"TransactionTest/internal/logger"

	"go.uber.org/zap"
)

// ITransactionService определяет интерфейс для работы с транзакциями.
// Используется для dependency injection в HTTP обработчиках.
type ITransactionService interface {
	// SendMoney отправляет деньги с одного кошелька на другой.
	// Возвращает код ошибки domain.ErrorCode.
	SendMoney(ctx context.Context, from, to string, amount float64) domain.ErrorCode

	// GetLastTransactions возвращает последние транзакции с ограничением по количеству.
	// Возвращает слайс транзакций и код ошибки.
	GetLastTransactions(ctx context.Context, limit int) ([]domain.Transaction, domain.ErrorCode)

	// GetTransactionById возвращает транзакцию по её ID.
	// Возвращает указатель на транзакцию и код ошибки.
	GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, domain.ErrorCode)

	// GetTransactionByInfo возвращает транзакцию по информации о ней.
	// Возвращает указатель на транзакцию и код ошибки.
	GetTransactionByInfo(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, domain.ErrorCode)

	// RemoveTransaction удаляет транзакцию по её ID.
	// Возвращает код ошибки.
	RemoveTransaction(ctx context.Context, id int64) domain.ErrorCode
}

// IWalletService определяет интерфейс для работы с кошельками.
type IWalletService interface {
	// CreateWallet создает новый кошелек с указанным балансом.
	// Возвращает адрес кошелька и код ошибки.
	CreateWallet(ctx context.Context, balance float64) (string, domain.ErrorCode)

	// GetBalance возвращает баланс кошелька по его адресу.
	// Возвращает баланс и код ошибки.
	GetBalance(ctx context.Context, address string) (float64, domain.ErrorCode)

	// GetWallet возвращает полную информацию о кошельке по его адресу.
	// Возвращает указатель на кошелек и код ошибки.
	GetWallet(ctx context.Context, address string) (*domain.Wallet, domain.ErrorCode)

	// UpdateBalance обновляет баланс кошелька.
	// Возвращает код ошибки.
	UpdateBalance(ctx context.Context, address string, newBalance float64) domain.ErrorCode

	// RemoveWallet удаляет кошелек по его адресу.
	// Возвращает код ошибки.
	RemoveWallet(ctx context.Context, address string) domain.ErrorCode
}

// Handler - HTTP обработчик для API.
// Содержит зависимости на сервисы транзакций и кошельков, а также логгер.
type Handler struct {
	transactionService ITransactionService
	walletService      IWalletService
	log                logger.Logger
}

// NewHandler создает новый экземпляр HTTP обработчика.
// Принимает сервисы транзакций и кошельков, а также логгер.
// Возвращает указатель на Handler.
func NewHandler(ts ITransactionService, ws IWalletService, l logger.Logger) *Handler {
	return &Handler{
		transactionService: ts,
		walletService:      ws,
		log:                l,
	}
}

// ErrorResponse представляет структуру ответа с ошибкой.
// Используется для стандартизации формата ошибок API.
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// writeJSON записывает JSON ответ в http.ResponseWriter.
// Устанавливает Content-Type header и кодирует данные в JSON.
// При ошибке кодирования логирует её.
func (h *Handler) writeJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error(ctx, "Failed to encode JSON response", zap.Error(err))
	}
}

// writeError записывает ошибку в формате JSON в http.ResponseWriter.
// Создает ErrorResponse с указанным статус кодом, кодом ошибки и сообщением.
func (h *Handler) writeError(ctx context.Context, w http.ResponseWriter, statusCode int, code domain.ErrorCode, message string) {
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Code:    string(code),
		Message: message,
	}
	h.writeJSON(ctx, w, statusCode, response)
}

// handleServiceError маппит коды ошибок домена в соответствующие HTTP статус коды и сообщения.
// Принимает код ошибки домена и название операции для логирования.
func (h *Handler) handleServiceError(ctx context.Context, w http.ResponseWriter, code domain.ErrorCode, operation string) {
	switch code {
	case domain.CodeOK:
		return
	case domain.CodeWalletNotFound:
		h.log.Warn(ctx, operation+": wallet not found")
		h.writeError(ctx, w, http.StatusNotFound, code, "Wallet not found")
	case domain.CodeTransactionNotFound:
		h.log.Warn(ctx, operation+": transaction not found")
		h.writeError(ctx, w, http.StatusNotFound, code, "Transaction not found")
	case domain.CodeInsufficientFunds:
		h.log.Warn(ctx, operation+": insufficient funds")
		h.writeError(ctx, w, http.StatusBadRequest, code, "Insufficient funds")
	case domain.CodeDuplicateWallet:
		h.log.Warn(ctx, operation+": duplicate wallet")
		h.writeError(ctx, w, http.StatusConflict, code, "Wallet already exists")
	case domain.CodeNegativeBalance:
		h.log.Warn(ctx, operation+": negative balance")
		h.writeError(ctx, w, http.StatusBadRequest, code, "Negative balance not allowed")
	case domain.CodeNegativeAmount:
		h.log.Warn(ctx, operation+": negative amount")
		h.writeError(ctx, w, http.StatusBadRequest, code, "Amount must be positive")
	case domain.CodeInvalidTransaction:
		h.log.Warn(ctx, operation+": invalid transaction")
		h.writeError(ctx, w, http.StatusBadRequest, code, "Invalid transaction")
	case domain.CodeInvalidLimit:
		h.log.Warn(ctx, operation+": invalid limit")
		h.writeError(ctx, w, http.StatusBadRequest, code, "Invalid limit parameter")
	case domain.CodeInternal:
		h.log.Error(ctx, operation+": internal error")
		h.writeError(ctx, w, http.StatusInternalServerError, code, "Internal server error")
	default:
		h.log.Error(ctx, operation+": unknown error code", zap.String("code", string(code)))
		h.writeError(ctx, w, http.StatusInternalServerError, domain.CodeInternal, "Internal server error")
	}
}
