package delivery

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "TransactionTest/internal/domain"
    "TransactionTest/internal/logger"
    "TransactionTest/internal/errors"

    "github.com/gorilla/mux"
    "go.uber.org/zap"
)

func (h *Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    appLogger := logger.GetLoggerFromCtx(ctx)

    var req struct {
        Balance float64 `json:"balance"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("CreateWalletHandler: failed to parse request body: %v", err),
        )
        return
    }   

    address, err := h.walletService.CreateWallet(ctx, req.Balance)
    if err != nil {
        http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("CreateWalletHandler: failed to create wallet: %v", err),
        )
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"address": address})
}

func (h *Handler) RemoveWallet(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    appLogger := logger.GetLoggerFromCtx(ctx)

    vars := mux.Vars(r)
    address, ok := vars["address"]
    if !ok {
        writeErrorResponse(w, ctx, errors.NewCustomError("Address is required", http.StatusBadRequest, nil), appLogger)
        return
    }

    err := h.walletService.RemoveWallet(ctx, address)
    if err != nil {
        http.Error(w, "Failed to remove wallet", http.StatusInternalServerError)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("RemoveWalletHandler: failed to remove wallet: %v", err),
        )
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
    ctx := h.GenerateRequestID(h.ctx)   

    vars := mux.Vars(r)
    address, ok := vars["address"]

    if !ok {
        http.Error(w, "Address is required", http.StatusBadRequest)
        return
    }

    var wallet *domain.Wallet

    wallet, err := h.walletService.GetWallet(ctx, address)
    if err != nil {
        http.Error(w, "Failed to get wallet", http.StatusInternalServerError)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("GetWalletInfoHandler: failed to get wallet: %v", err),
        )
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(wallet)
}

func (h *Handler) GetTransactionById(w http.ResponseWriter, r *http.Request) {
    ctx := h.GenerateRequestID(h.ctx)    

    vars := mux.Vars(r)
    idStr, ok := vars["id"]
    if !ok {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var transaction *domain.Transaction

    transaction, err = h.transactionService.GetTransactionById(ctx, id)
    if err != nil {
        http.Error(w, "Transaction not found", http.StatusNotFound)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("GetTransactionByIdHandler: transaction not found: %v", err),
        )
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transaction)
}

func (h *Handler) RemoveTransaction(w http.ResponseWriter, r *http.Request) {
    ctx := h.GenerateRequestID(h.ctx)   

    vars := mux.Vars(r)
    idStr, ok := vars["id"]
    if !ok {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    err = h.transactionService.RemoveTransaction(ctx, id)
    if err != nil {
        http.Error(w, "Failed to remove transaction", http.StatusInternalServerError)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("RemoveTransactionHandler: failed to remove transaction: %v", err),
        )
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetTransactionByInfo(w http.ResponseWriter, r *http.Request) {
    ctx := h.GenerateRequestID(h.ctx)   

    vars := mux.Vars(r)
    from, ok := vars["from"]
    if !ok {
        http.Error(w, "Sender address (from) is required", http.StatusBadRequest)
        return
    }

    to, ok := vars["to"]
    if !ok {
        http.Error(w, "Receiver address (to) is required", http.StatusBadRequest)
        return
    }

    createdAtStr, ok := vars["createdAt"]
    if !ok {
        http.Error(w, "Transaction timestamp (createdAt) is required", http.StatusBadRequest)
        return
    }

    createdAt, err := time.Parse(time.RFC3339, createdAtStr)
    if err != nil {
        http.Error(w, "Invalid timestamp format. Use RFC3339 format (e.g., 2024-02-10T15:04:05Z)", http.StatusBadRequest)
        return
    }

    transaction, err := h.transactionService.GetTransactionByInfo(ctx, from, to, createdAt)
    if err != nil {
        http.Error(w, "Transaction not found", http.StatusNotFound)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("GetTransactionByInfoHandler: transaction not found: %v", err),
        )
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transaction)
}