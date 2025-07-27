package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"TransactionTest/internal/domain"
	"TransactionTest/internal/delivery/dto"
	"TransactionTest/internal/logger"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) SendMoney(ctx context.Context, from, to string, amount float64) domain.ErrorCode {
	args := m.Called(ctx, from, to, amount)
	return args.Get(0).(domain.ErrorCode)
}

func (m *MockTransactionService) GetLastTransactions(ctx context.Context, limit int) ([]domain.Transaction, domain.ErrorCode) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]domain.Transaction), args.Get(1).(domain.ErrorCode)
}

func (m *MockTransactionService) GetTransactionById(ctx context.Context, id int64) (*domain.Transaction, domain.ErrorCode) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Transaction), args.Get(1).(domain.ErrorCode)
}

func (m *MockTransactionService) GetTransactionByInfo(ctx context.Context, from, to string, createdAt time.Time) (*domain.Transaction, domain.ErrorCode) {
	args := m.Called(ctx, from, to, createdAt)
	return args.Get(0).(*domain.Transaction), args.Get(1).(domain.ErrorCode)
}

func (m *MockTransactionService) RemoveTransaction(ctx context.Context, id int64) domain.ErrorCode {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.ErrorCode)
}

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) CreateWallet(ctx context.Context, balance float64) (string, domain.ErrorCode) {
	args := m.Called(ctx, balance)
	return args.String(0), args.Get(1).(domain.ErrorCode)
}

func (m *MockWalletService) GetBalance(ctx context.Context, address string) (float64, domain.ErrorCode) {
	args := m.Called(ctx, address)
	return args.Get(0).(float64), args.Get(1).(domain.ErrorCode)
}

func (m *MockWalletService) GetWallet(ctx context.Context, address string) (*domain.Wallet, domain.ErrorCode) {
	args := m.Called(ctx, address)
	return args.Get(0).(*domain.Wallet), args.Get(1).(domain.ErrorCode)
}

func (m *MockWalletService) UpdateBalance(ctx context.Context, address string, newBalance float64) domain.ErrorCode {
	args := m.Called(ctx, address, newBalance)
	return args.Get(0).(domain.ErrorCode)
}

func (m *MockWalletService) RemoveWallet(ctx context.Context, address string) domain.ErrorCode {
	args := m.Called(ctx, address)
	return args.Get(0).(domain.ErrorCode)
}

func setupTest() (*Handler, *MockTransactionService, *MockWalletService, logger.Logger) {
	mockTS := new(MockTransactionService)
	mockWS := new(MockWalletService)
	
	// Create test logger
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel) // Only errors in tests
	log, _ := logger.New(&config)
	
	handler := NewHandler(mockTS, mockWS, log)
	
	return handler, mockTS, mockWS, log
}

func TestCreateWallet_Success(t *testing.T) {
	handler, _, mockWS, _ := setupTest()
	
	req := dto.CreateWalletRequest{
		Balance: 100.0,
	}
	
	expectedAddress := "550e8400-e29b-41d4-a716-446655440000"
	mockWS.On("CreateWallet", mock.Anything, req.Balance).Return(expectedAddress, domain.CodeOK)
	
	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/api/wallet/create", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	
	recorder := httptest.NewRecorder()
	handler.CreateWallet(recorder, request)
	
	assert.Equal(t, http.StatusCreated, recorder.Code)
	
	var response dto.CreateWalletResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, expectedAddress, response.Address)
	
	mockWS.AssertExpectations(t)
}

func TestCreateWallet_ValidationError(t *testing.T) {
	handler, _, _, _ := setupTest()
	
	req := dto.CreateWalletRequest{
		Balance: -100.0,
	}
	
	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/api/wallet/create", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	
	recorder := httptest.NewRecorder()
	handler.CreateWallet(recorder, request)
	
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	
	var response ErrorResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, string(domain.CodeInvalidRequestBody), response.Code)
}

func TestSendMoney_Success(t *testing.T) {
	handler, mockTS, _, _ := setupTest()
	
	req := dto.SendMoneyRequest{
		From:   "550e8400-e29b-41d4-a716-446655440000",
		To:     "550e8400-e29b-41d4-a716-446655440001",
		Amount: 100.0,
	}
	
	mockTS.On("SendMoney", mock.Anything, req.From, req.To, req.Amount).Return(domain.CodeOK)
	
	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/api/send", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	
	recorder := httptest.NewRecorder()
	handler.SendMoney(recorder, request)
	
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "Money sent successfully", response["message"])
	
	mockTS.AssertExpectations(t)
}

func TestSendMoney_ValidationError(t *testing.T) {
	handler, _, _, _ := setupTest()
	
	req := dto.SendMoneyRequest{
		From:   "invalid-uuid",
		To:     "550e8400-e29b-41d4-a716-446655440001",
		Amount: 100.0,
	}
	
	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/api/send", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	
	recorder := httptest.NewRecorder()
	handler.SendMoney(recorder, request)
	
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	
	var response ErrorResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, string(domain.CodeInvalidRequestBody), response.Code)
}

func TestGetBalance_Success(t *testing.T) {
	handler, _, mockWS, _ := setupTest()
	
	address := "550e8400-e29b-41d4-a716-446655440000"
	expectedBalance := 100.0
	
	mockWS.On("GetBalance", mock.Anything, address).Return(expectedBalance, domain.CodeOK)
	
	request := httptest.NewRequest("GET", "/api/wallet/"+address+"/balance", nil)
	vars := map[string]string{"address": address}
	request = mux.SetURLVars(request, vars)
	
	recorder := httptest.NewRecorder()
	handler.GetBalance(recorder, request)
	
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response dto.BalanceResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, expectedBalance, response.Balance)
	
	mockWS.AssertExpectations(t)
}

func TestGetBalance_WalletNotFound(t *testing.T) {
	handler, _, mockWS, _ := setupTest()
	
	address := "550e8400-e29b-41d4-a716-446655440000"
	
	mockWS.On("GetBalance", mock.Anything, address).Return(0.0, domain.CodeWalletNotFound)
	
	request := httptest.NewRequest("GET", "/api/wallet/"+address+"/balance", nil)
	vars := map[string]string{"address": address}
	request = mux.SetURLVars(request, vars)
	
	recorder := httptest.NewRecorder()
	handler.GetBalance(recorder, request)
	
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	
	var response ErrorResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, string(domain.CodeWalletNotFound), response.Code)
	
	mockWS.AssertExpectations(t)
}

func TestGetLastTransactions_Success(t *testing.T) {
	handler, mockTS, _, _ := setupTest()
	
	limit := 10
	expectedTransactions := []domain.Transaction{
		{
			Id:        1,
			From:      "550e8400-e29b-41d4-a716-446655440000",
			To:        "550e8400-e29b-41d4-a716-446655440001",
			Amount:    100.0,
			CreatedAt: time.Now(),
		},
	}
	
	mockTS.On("GetLastTransactions", mock.Anything, limit).Return(expectedTransactions, domain.CodeOK)
	
	request := httptest.NewRequest("GET", "/api/transactions?count=10", nil)
	
	recorder := httptest.NewRecorder()
	handler.GetLastTransactions(recorder, request)
	
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response []dto.TransactionResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Len(t, response, 1)
	assert.Equal(t, expectedTransactions[0].Id, response[0].Id)
	
	mockTS.AssertExpectations(t)
}

func TestGetLastTransactions_InvalidLimit(t *testing.T) {
	handler, _, _, _ := setupTest()
	
	request := httptest.NewRequest("GET", "/api/transactions?count=0", nil)
	
	recorder := httptest.NewRecorder()
	handler.GetLastTransactions(recorder, request)
	
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	
	var response ErrorResponse
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, string(domain.CodeInvalidLimit), response.Code)
} 