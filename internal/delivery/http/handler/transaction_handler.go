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

func (h *Handler) SendMoney(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "SendMoney: "

	var req dto.SendMoneyRequest
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

	svcCode := h.transactionService.SendMoney(ctx, req.From, req.To, req.Amount)
	if svcCode != domain.CodeOK {
		h.log.Warn(
			ctx,
			op+"service returned error",
			zap.String("error_code", string(svcCode)),
		)
		h.handleServiceError(ctx, w, svcCode, "SendMoney")
		return
	}

	h.log.Info(
		ctx,
		op+"transaction completed successfully",
		zap.String("from", req.From),
		zap.String("to", req.To),
		zap.Float64("amount", req.Amount),
	)
	h.writeJSON(ctx, w, http.StatusOK, map[string]string{"message": "Money sent successfully"})
}

func (h *Handler) GetLastTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "GetLastTransactions: "

	count, code, msg := h.parseAndValidateCount(ctx, r, op)
    if code != 0 {
        h.writeError(ctx, w, code, domain.CodeInvalidRequestBody, msg)
        return
    }

	h.log.Info(
		ctx,
		op+"received request",
		zap.Int("count", count),
	)

	transactions, svcCode := h.transactionService.GetLastTransactions(ctx, count)
	if svcCode != domain.CodeOK {
		h.log.Warn(
			ctx,
			op+"service returned error",
			zap.String("error_code", string(svcCode)),
		)
		h.handleServiceError(ctx, w, svcCode, "GetLastTransactions")
		return
	}

	response := make([]dto.TransactionResponse, len(transactions))
	for i, t := range transactions {
		response[i] = dto.TransactionResponse{
			Id:        t.Id,
			From:      t.From,
			To:        t.To,
			Amount:    t.Amount,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
		}
	}

	h.log.Info(
		ctx,
		op+"transactions retrieved successfully",
		zap.Int("count", len(transactions)),
	)
	h.writeJSON(ctx, w, http.StatusOK, response)
}

func (h *Handler) GetTransactionById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "GetTransactionById: "

	id, code, msg := h.parseAndValidateID(ctx, r, op)
    if code != 0 {
        h.writeError(ctx, w, code, domain.CodeInvalidRequestBody, msg)
        return
    }

	h.log.Info(
		ctx,
		op+"received request",
		zap.Int64("id", id),
	)

	transaction, svcCode := h.transactionService.GetTransactionById(ctx, id)
	if svcCode != domain.CodeOK {
		h.log.Warn(
			ctx,
			op+"service returned error",
			zap.String("error_code", string(svcCode)),
		)
		h.handleServiceError(ctx, w, svcCode, "GetTransactionById")
		return
	}

	response := dto.TransactionResponse{
		Id:        transaction.Id,
		From:      transaction.From,
		To:        transaction.To,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt.Format(time.RFC3339),
	}

	h.log.Info(
		ctx,
		op+"transaction retrieved successfully",
		zap.Int64("id", id),
	)
	h.writeJSON(ctx, w, http.StatusOK, response)
}

func (h *Handler) GetTransactionByInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "GetTransactionByInfo: "

	var req dto.GetTransactionByInfoRequest
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
		zap.String("from", req.From),
		zap.String("to", req.To),
		zap.String("createdAt", req.CreatedAt),
	)

	createdAt, err := time.Parse(time.RFC3339, req.CreatedAt)
	if err != nil {
		h.log.Warn(
			ctx,
			op+"invalid created_at format",
			zap.String("createdAt", req.CreatedAt),
			zap.Error(err),
		)
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid created_at format")
		return
	}

	transaction, srvCode := h.transactionService.GetTransactionByInfo(ctx, req.From, req.To, createdAt)
	if srvCode != domain.CodeOK {
		h.log.Warn(
			ctx,
			op+"service returned error",
			zap.String("error_code", string(srvCode)),
		)
		h.handleServiceError(ctx, w, srvCode, "GetTransactionByInfo")
		return
	}

	response := dto.TransactionResponse{
		Id:        transaction.Id,
		From:      transaction.From,
		To:        transaction.To,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt.Format(time.RFC3339),
	}

	h.log.Info(
		ctx,
		op+"transaction retrieved successfully",
		zap.String("from", req.From),
		zap.String("to", req.To),
		zap.String("createdAt", req.CreatedAt),
	)
	h.writeJSON(ctx, w, http.StatusOK, response)
}

func (h *Handler) RemoveTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	const op = "RemoveTransaction: "

	id, code, msg := h.parseAndValidateID(ctx, r, op)
    if code != 0 {
        h.writeError(ctx, w, code, domain.CodeInvalidRequestBody, msg)
        return
    }

	h.log.Info(
		ctx,
		op+"received request",
		zap.Int64("id", id),
	)

	srvCode := h.transactionService.RemoveTransaction(ctx, id)
	if srvCode != domain.CodeOK {
		h.log.Warn(
			ctx,
			op+"service returned error",
			zap.String("error_code", string(srvCode)),
		)
		h.handleServiceError(ctx, w, srvCode, "RemoveTransaction")
		return
	}

	h.log.Info(
		ctx,
		op+"transaction removed successfully",
		zap.Int64("id", id),
	)
	h.writeJSON(ctx, w, http.StatusOK, map[string]string{"message": "Transaction removed successfully"})
} 