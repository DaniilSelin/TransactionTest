package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"TransactionTest/internal/delivery/dto"
	"TransactionTest/internal/domain"
	"TransactionTest/internal/delivery/http"
	"github.com/stretchr/testify/assert"
)

type mockTransactionService struct {
	SendMoneyFunc                func(ctx interface{}, from, to string, amount float64) domain.ErrorCode
	GetLastTransactionsFunc      func(ctx interface{}, limit int) ([]domain.Transaction, domain.ErrorCode)
	GetTransactionByIdFunc       func(ctx interface{}, id int64) (*domain.Transaction, domain.ErrorCode)
	GetTransactionByInfoFunc     func(ctx interface{}, from, to string, amount float64, createdAt time.Time) (*domain.Transaction, domain.ErrorCode)
	RemoveTransactionFunc        func(ctx interface{}, id int64) domain.ErrorCode
}

// Реализация методов mockTransactionService ... (опущено для краткости)

func TestTransactionHandler_SendMoney_Success(t *testing.T) {
	h := http.Handler{ /* инициализация с моками */ }
	body, _ := json.Marshal(dto.SendMoneyRequest{From: "uuid1", To: "uuid2", Amount: 10})
	req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.SendMoney(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Аналогично: ошибки декодирования, ошибки валидации, ошибки сервиса, edge-cases
// Аналогично для GetLastTransactions, GetTransactionById, GetTransactionByInfo, RemoveTransaction
