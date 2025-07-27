package handler

import (
	"context"
	"net/http"
	"strconv"

	"TransactionTest/internal/delivery/dto"
	"TransactionTest/internal/delivery/validator"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// parseAndValidateAddress извлекает извлекает {address} из URL, оборачивает в DTO и валидирует.
// При ошибке возвращает HTTP‑код и сообщение, чтобы handler мог сразу ответить.
func (h *Handler) parseAndValidateAddress(
	ctx context.Context,
	r *http.Request,
	operation string,
) (string, int, string) {
	vars := mux.Vars(r)
	addr, ok := vars["address"]
	if !ok {
		h.log.Warn(ctx, operation+"address not provided")
		return "", http.StatusBadRequest, "address not provided"
	}
	p := dto.AddressPath{Address: addr}
	if err := validator.ValidateStruct(p); err != nil {
		h.log.Warn(ctx, operation+"path validation failed", zap.Error(err))
		return "", http.StatusBadRequest, err.Error()
	}
	return addr, 0, ""
}

// parseAndValidateID извлекает параметр ?id из URL, оборачивает в DTO и валидирует.
// При ошибке возвращает HTTP‑код и сообщение, чтобы handler мог сразу ответить.
func (h *Handler) parseAndValidateID(
	ctx context.Context,
	r *http.Request,
	operation string,
) (int64, int, string) {
	vars := mux.Vars(r)
	s, ok := vars["id"]
	if !ok {
		h.log.Warn(ctx, operation+"id not provided")
		return 0, http.StatusBadRequest, "id not provided"
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		h.log.Warn(ctx, operation+"id parse failed", zap.Error(err))
		return 0, http.StatusBadRequest, "invalid id"
	}
	p := dto.TransactionID{ID: id}
	if err := validator.ValidateStruct(p); err != nil {
		h.log.Warn(ctx, operation+"path validation failed", zap.Error(err))
		return 0, http.StatusBadRequest, err.Error()
	}
	return id, 0, ""
}

// parseAndValidateCount извлекает параметр ?count из URL, оборачивает в DTO и валидирует.
// При ошибке возвращает HTTP‑код и сообщение, чтобы handler мог сразу ответить.
func (h *Handler) parseAndValidateCount(
	ctx context.Context,
	r *http.Request,
	operation string,
) (int, int, string) {
	raw := r.URL.Query().Get("count")
	if raw == "" {
		h.log.Warn(ctx, operation+"count not provided")
		return 0, http.StatusBadRequest, "count is required"
	}
	count, err := strconv.Atoi(raw)
	if err != nil {
		h.log.Warn(ctx, operation+"count parse failed", zap.Error(err))
		return 0, http.StatusBadRequest, "invalid count"
	}
	p := dto.CountQuery{Count: count}
	if err := validator.ValidateStruct(p); err != nil {
		h.log.Warn(ctx, operation+"count validation failed", zap.Any("errors", err))
		return 0, http.StatusBadRequest, err.Error()
	}
	return count, 0, ""
}
