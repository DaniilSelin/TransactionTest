package handler

import (
    "encoding/json"
    "net/http"
    "TransactionTest/internal/domain"
)

type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code,omitempty"`
    Message string `json:"message"`
}

func (h *Handler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    _ = json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, statusCode int, code domain.ErrorCode, message string) {
    resp := ErrorResponse{
        Error:   http.StatusText(statusCode),
        Code:    string(code),
        Message: message,
    }
    h.writeJSON(w, statusCode, resp)
}
