package http

import (
	"TransactionTest/internal/logger"
	"github.com/gorilla/mux"
	httpBase "net/http"
)

type IHanlder interface {
	SendMoney(w httpBase.ResponseWriter, r *httpBase.Request)
	GetLastTransactions(w httpBase.ResponseWriter, r *httpBase.Request)
	GetBalance(w httpBase.ResponseWriter, r *httpBase.Request)

	GetTransactionById(w httpBase.ResponseWriter, r *httpBase.Request)
	RemoveTransaction(w httpBase.ResponseWriter, r *httpBase.Request)
	GetTransactionByInfo(w httpBase.ResponseWriter, r *httpBase.Request)

	CreateWallet(w httpBase.ResponseWriter, r *httpBase.Request)
	GetWallet(w httpBase.ResponseWriter, r *httpBase.Request)
	RemoveWallet(w httpBase.ResponseWriter, r *httpBase.Request)
	UpdateBalance(w httpBase.ResponseWriter, r *httpBase.Request)
}

func NewRouter(h IHanlder, log logger.Logger) *mux.Router {
	r := mux.NewRouter()

	// Добавляем middleware
	r.Use(RequestIDMiddleware)
	r.Use(LoggingMiddleware(log))
	r.Use(RecoveryMiddleware(log))

	api := r.PathPrefix("/api").Subrouter()

	// Пути указанные в ТЗ
	api.HandleFunc("/send", h.SendMoney).Methods(httpBase.MethodPost)
	api.HandleFunc("/transactions", h.GetLastTransactions).Methods(httpBase.MethodGet).Queries("count", "{count}")
	api.HandleFunc("/wallet/{address}/balance", h.GetBalance).Methods(httpBase.MethodGet)

	// Дополнительные пути
	api.HandleFunc("/transaction/{id}", h.GetTransactionById).Methods(httpBase.MethodGet)
	api.HandleFunc("/transaction/{id}", h.RemoveTransaction).Methods(httpBase.MethodDelete)
	api.HandleFunc("/transaction/{from}/{to}/{createdAt}", h.GetTransactionByInfo).Methods(httpBase.MethodGet)

	// Ожидает на вход - { "balance": x.x }
	api.HandleFunc("/wallet/create", h.CreateWallet).Methods(httpBase.MethodPost)

	api.HandleFunc("/wallet/{address}", h.GetWallet).Methods(httpBase.MethodGet)
	api.HandleFunc("/wallet/{address}", h.RemoveWallet).Methods(httpBase.MethodDelete)
	api.HandleFunc("/wallet/{address}/balance", h.UpdateBalance).Methods(httpBase.MethodPut)

	return r
}
