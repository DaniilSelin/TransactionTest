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

type mockWalletService struct {
	CreateWalletFunc            func(ctx interface{}, balance float64) (string, domain.ErrorCode)
	GetBalanceFunc              func(ctx interface{}, address string) (float64, domain.ErrorCode)
	GetWalletFunc               func(ctx interface{}, address string) (*domain.Wallet, domain.ErrorCode)
	UpdateBalanceFunc           func(ctx interface{}, address string, balance float64) domain.ErrorCode
	RemoveWalletFunc            func(ctx interface{}, address string) domain.ErrorCode
}

// Реализация методов mockWalletService ... (опущено для краткости)

func TestWalletHandler_CreateWallet_Success(t *testing.T) {
	h := http.Handler{ /* инициализация с моками */ }
	body, _ := json.Marshal(dto.CreateWalletRequest{Balance: 100})
	req := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateWallet(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

// Аналогично: ошибки декодирования, ошибки валидации, ошибки сервиса, edge-cases
// Аналогично для GetBalance, GetWallet, UpdateBalance, RemoveWallet
