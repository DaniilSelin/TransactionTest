package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"TransactionTest/internal/delivery/dto"
	"TransactionTest/internal/delivery/validator"
	"TransactionTest/internal/domain"

	"go.uber.org/zap"
)

// CreateWallet обрабатывает HTTP POST запрос для создания нового кошелька.
//
// Принимает JSON в теле запроса:
//
//	{
//	  "balance": 100.50
//	}
//
// Возможные коды ответа:
//   - 201 Created: кошелек успешно создан
//   - 400 Bad Request: ошибка валидации (отрицательный баланс)
//   - 409 Conflict: кошелек уже существует (невозможно)
//   - 500 Internal Server Error: внутренняя ошибка сервера
//
// Пример успешного ответа:
//
//	{
//	  "address": "550e8400-e29b-41d4-a716-446655440000"
//	}
func (h *Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "CreateWallet: "

	var req dto.CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(
			ctx,
			op+"failed to decode JSON",
			zap.Error(err),
		)
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid JSON")
		return
	}

	h.log.Info(
		ctx,
		op+"received request",
		zap.Any("payload", req),
	)

	if err := validator.ValidateStruct(req); err != nil {
		h.log.Warn(
			ctx,
			op+"validation failed",
			zap.Any("errors", err),
		)
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, err.Error())
		return
	}

	address, svcCode := h.walletService.CreateWallet(ctx, req.Balance)
	if svcCode != domain.CodeOK {
		h.log.Warn(
			ctx,
			op+"service returned error",
			zap.String("error_code", string(svcCode)),
		)
		h.handleServiceError(ctx, w, svcCode, "CreateWallet")
		return
	}

	h.log.Info(
		ctx,
		op+"wallet created successfully",
		zap.String("address", address),
		zap.Float64("balance", req.Balance),
	)
	h.writeJSON(ctx, w, http.StatusCreated, dto.CreateWalletResponse{Address: address})
}

// GetBalance обрабатывает HTTP GET запрос для получения баланса кошелька.
//
// Path параметры:
//   - address: адрес кошелька (UUID4, обязательный)
//
// URL: GET /api/wallet/uuid4/balance
//
// Возможные коды ответа:
//   - 200 OK: баланс успешно получен
//   - 400 Bad Request: неверный адрес
//   - 404 Not Found: кошелек не найден
//   - 500 Internal Server Error: внутренняя ошибка сервера
//
// Пример успешного ответа:
//
//	{
//	  "balance": 100.50
//	}
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "GetBalance: "

	address, code, msg := h.parseAndValidateAddress(ctx, r, op)
	if code != 0 {
		h.writeError(ctx, w, code, domain.CodeInvalidRequestBody, msg)
		return
	}

	h.log.Info(
		ctx,
		op+"received request",
		zap.String("address", address),
	)

	balance, svcCode := h.walletService.GetBalance(ctx, address)
	if svcCode != domain.CodeOK {
		h.log.Warn(
			ctx,
			op+"service returned error",
			zap.String("error_code", string(svcCode)),
		)
		h.handleServiceError(ctx, w, svcCode, "GetBalance")
		return
	}

	h.log.Info(
		ctx,
		op+"balance retrieved successfully",
		zap.String("address", address),
		zap.Float64("balance", balance),
	)
	h.writeJSON(ctx, w, http.StatusOK, dto.BalanceResponse{Balance: balance})
}

// GetWallet обрабатывает HTTP GET запрос для получения полной информации о кошельке.
//
// Path параметры:
//   - address: адрес кошелька (UUID, обязательный)
//
// URL: GET /api/wallet/uuid4
//
// Возможные коды ответа:
//   - 200 OK: информация о кошельке получена
//   - 400 Bad Request: неверный адрес
//   - 404 Not Found: кошелек не найден
//   - 500 Internal Server Error: внутренняя ошибка сервера
//
// Пример успешного ответа:
//
//	{
//	  "address": "550e8400-e29b-41d4-a716-446655440000",
//	  "balance": 100.50,
//	  "created_at": "2023-01-01T12:00:00Z"
//	}
func (h *Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "GetWallet: "

	address, code, msg := h.parseAndValidateAddress(ctx, r, op)
	if code != 0 {
		h.writeError(ctx, w, code, domain.CodeInvalidRequestBody, msg)
		return
	}

	h.log.Info(
		ctx,
		op+"received request",
		zap.String("address", address),
	)

	wallet, svcCode := h.walletService.GetWallet(ctx, address)
	if svcCode != domain.CodeOK {
		h.log.Warn(
			ctx,
			op+"service returned error",
			zap.String("error_code", string(svcCode)),
		)
		h.handleServiceError(ctx, w, svcCode, "GetWallet")
		return
	}

	response := dto.WalletResponse{
		Address:   wallet.Address,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt.Format(time.RFC3339),
	}

	h.log.Info(
		ctx,
		op+"wallet retrieved successfully",
		zap.String("address", address),
	)
	h.writeJSON(ctx, w, http.StatusOK, response)
}

// UpdateBalance обрабатывает HTTP PUT запрос для обновления баланса кошелька.
//
// Path параметры:
//   - address: адрес кошелька (UUID, обязательный)
//
// Принимает JSON в теле запроса:
//
//	{
//	  "balance": 200.75
//	}
//
// URL: PUT /api/wallet/550e8400-e29b-41d4-a716-446655440000/balance
//
// Возможные коды ответа:
//   - 200 OK: баланс успешно обновлен
//   - 400 Bad Request: ошибка валидации или отрицательный баланс
//   - 404 Not Found: кошелек не найден
//   - 500 Internal Server Error: внутренняя ошибка сервера
//
// Пример успешного ответа:
//
//	{
//	  "message": "Balance updated successfully"
//	}
func (h *Handler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "UpdateBalance: "

	address, code, msg := h.parseAndValidateAddress(ctx, r, op)
	if code != 0 {
		h.writeError(ctx, w, code, domain.CodeInvalidRequestBody, msg)
		return
	}

	var req dto.UpdateBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(
			ctx,
			op+"failed to decode JSON",
			zap.Error(err),
		)
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid JSON")
		return
	}

	h.log.Info(
		ctx,
		op+"received request",
		zap.String("address", address),
		zap.Any("payload", req),
	)

	if err := validator.ValidateStruct(req); err != nil {
		h.log.Warn(
			ctx,
			op+"request validation failed",
			zap.Any("errors", err),
		)
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, err.Error())
		return
	}

	svcCode := h.walletService.UpdateBalance(ctx, address, req.Balance)
	if svcCode != domain.CodeOK {
		h.log.Warn(
			ctx,
			op+"service returned error",
			zap.String("error_code", string(svcCode)),
		)
		h.handleServiceError(ctx, w, svcCode, "UpdateBalance")
		return
	}

	h.log.Info(
		ctx,
		op+"balance updated successfully",
		zap.String("address", address),
		zap.Float64("new_balance", req.Balance),
	)
	h.writeJSON(ctx, w, http.StatusOK, map[string]string{"message": "Balance updated successfully"})
}

// RemoveWallet обрабатывает HTTP DELETE запрос для удаления кошелька.
//
// Path параметры:
//   - address: адрес кошелька (UUID, обязательный)
//
// URL: DELETE /api/wallet/550e8400-e29b-41d4-a716-446655440000
//
// Возможные коды ответа:
//   - 200 OK: кошелек успешно удален
//   - 400 Bad Request: неверный адрес
//   - 404 Not Found: кошелек не найден
//   - 500 Internal Server Error: внутренняя ошибка сервера
//
// Пример успешного ответа:
//
//	{
//	  "message": "Wallet removed successfully"
//	}
func (h *Handler) RemoveWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "RemoveWallet: "

	address, code, msg := h.parseAndValidateAddress(ctx, r, op)
	if code != 0 {
		h.writeError(ctx, w, code, domain.CodeInvalidRequestBody, msg)
		return
	}

	h.log.Info(
		ctx,
		op+"received request",
		zap.String("address", address),
	)

	svcCode := h.walletService.RemoveWallet(ctx, address)
	if svcCode != domain.CodeOK {
		h.log.Warn(
			ctx,
			op+"service returned error",
			zap.String("error_code", string(svcCode)),
		)
		h.handleServiceError(ctx, w, svcCode, "RemoveWallet")
		return
	}

	h.log.Info(
		ctx,
		op+"wallet removed successfully",
		zap.String("address", address),
	)
	h.writeJSON(ctx, w, http.StatusOK, map[string]string{"message": "Wallet removed successfully"})
}
