package delivery

import (
    "encoding/json"
    "net/http"
    "strconv"
    "context"
    "fmt"

    "TransactionTest/internal/service"
    "TransactionTest/internal/logger"
    "TransactionTest/internal/domain"

    "github.com/gorilla/mux"
    "github.com/google/uuid"
)

type Handler struct {
    transactionService *service.TransactionService
    walletService      *service.WalletService
    ctx                context.Context
}

func NewHandler(ctx context.Context, ts *service.TransactionService, ws *service.WalletService) *Handler {
    return &Handler{
        transactionService: ts,
        walletService:      ws,
        ctx: ctx,
    }
}

func (h *Handler) GenerateRequestID(ctx context.Context) context.Context {
    ctx = context.WithValue(ctx, logger.RequestID, uuid.New().String())
    return ctx
} 

func (h *Handler) SendMoney(w http.ResponseWriter, r *http.Request) {
    ctx := h.GenerateRequestID(h.ctx)  

    var req struct {
        From   string  `json:"from"`
        To     string  `json:"to"`
        Amount float64 `json:"amount"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("SendMoneyHandler: failed to decode request: %v", err),
        )
        return
    }

    err := h.transactionService.SendMoney(ctx, req.From, req.To, req.Amount)
    if err != nil {
        http.Error(w, "Transaction failed", http.StatusBadRequest)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("SendMoneyHandler: transaction failed: %v", err),
        )
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetLastTransactions(w http.ResponseWriter, r *http.Request) {
    ctx := h.GenerateRequestID(h.ctx)    

    countStr := r.URL.Query().Get("count")
    count, err := strconv.Atoi(countStr)
    if err != nil || count <= 0 {
        http.Error(w, "Invalid count parameter", http.StatusBadRequest)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("GetLastTransactionsHandler: invalid count parameter: %v", err),
        )
        return
    }

    transactions, err := h.transactionService.GetLastTransactions(ctx, count)
    if err != nil {
        http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("GetLastTransactionsHandler: failed to fetch transactions: %v", err),
        )
        return
    }

    logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("GetLastTransactionsHandler: successfule get LastTranhsaction"),
        )

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
    ctx := h.GenerateRequestID(h.ctx)    

    vars := mux.Vars(r)
    address, ok := vars["address"]
    if !ok {
        http.Error(w, "Wallet address is required", http.StatusBadRequest)
        return
    }

    balance, err := h.walletService.GetBalance(ctx, address)
    if err != nil {
        http.Error(w, "Failed to fetch balance", http.StatusInternalServerError)
        logger.GetLoggerFromCtx(ctx).Info(ctx, 
            fmt.Sprintf("GetBalanceHandler: failed to fetch balance: %v", err),
        )
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}