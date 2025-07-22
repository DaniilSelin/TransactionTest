package delivery

import (
	"net/http"
	"context"

	"github.com/gorilla/mux"
	"TransactionTest/internal/service"
)

func NewRouter(ctx context.Context, transactionService *service.TransactionService, walletService *service.WalletService) *mux.Router {
	r := mux.NewRouter()
	h := NewHandler(ctx, transactionService, walletService)

	api := r.PathPrefix("/api").Subrouter()

	// Пути указанные в ТЗ
	api.HandleFunc("/send", h.SendMoney).Methods(http.MethodPost)
	api.HandleFunc("/transactions", h.GetLastTransactions).Methods(http.MethodGet).Queries("count", "{count}")
	api.HandleFunc("/wallet/{address}/balance", h.GetBalance).Methods(http.MethodGet)

	// Дополнительные пути, необходимые 
	api.HandleFunc("/transaction/{id}", h.GetTransactionById).Methods(http.MethodGet)
	api.HandleFunc("/transaction/{id}", h.RemoveTransaction).Methods(http.MethodDelete)
	api.HandleFunc("/transaction/{from}/{to}/{createdAt}", h.GetTransactionByInfo).Methods(http.MethodGet)
	
	// Ожидает на вход - { "balance": x.x }
	api.HandleFunc("/wallet/create", h.CreateWallet).Methods(http.MethodPost)

	api.HandleFunc("/wallet/get/{address}", h.GetWallet).Methods(http.MethodGet)
	api.HandleFunc("/wallet/remove/{address}", h.RemoveWallet).Methods(http.MethodDelete)

	lmw := NewMiddlerWare()
	r.Use(lmw.loggerMiddleware)

	return r
}
