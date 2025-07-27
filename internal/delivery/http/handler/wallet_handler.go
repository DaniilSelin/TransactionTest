package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"TransactionTest/internal/domain"
	"TransactionTest/internal/delivery/dto"
	"TransactionTest/internal/delivery/validator"

	"go.uber.org/zap"
)

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