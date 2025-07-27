package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"TransactionTest/internal/logger"
)

type IHanlder interface {
	SendMoney(w http.ResponseWriter, r *http.Request)
	GetLastTransactions(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)

	GetTransactionById(w http.ResponseWriter, r *http.Request)
	RemoveTransaction(w http.ResponseWriter, r *http.Request)
	GetTransactionByInfo(w http.ResponseWriter, r *http.Request)

	CreateWallet(w http.ResponseWriter, r *http.Request)
	GetWallet(w http.ResponseWriter, r *http.Request)
	RemoveWallet(w http.ResponseWriter, r *http.Request)
	UpdateBalance(w http.ResponseWriter, r *http.Request)
}

func NewRouter(h IHanlder, log logger.Logger) *mux.Router {
	r := mux.NewRouter()
	
	// Добавляем middleware
	r.Use(RequestIDMiddleware)
	r.Use(LoggingMiddleware(log))
	r.Use(RecoveryMiddleware(log))
	
	api := r.PathPrefix("/api").Subrouter()

	// Пути указанные в ТЗ
	api.HandleFunc("/send", h.SendMoney).Methods(http.MethodPost)
	api.HandleFunc("/transactions", h.GetLastTransactions).Methods(http.MethodGet).Queries("count", "{count}")
	api.HandleFunc("/wallet/{address}/balance", h.GetBalance).Methods(http.MethodGet)

	// Дополнительные пути 
	api.HandleFunc("/transaction/{id}", h.GetTransactionById).Methods(http.MethodGet)
	api.HandleFunc("/transaction/{id}", h.RemoveTransaction).Methods(http.MethodDelete)
	api.HandleFunc("/transaction/{from}/{to}/{createdAt}", h.GetTransactionByInfo).Methods(http.MethodGet)
	
	// Ожидает на вход - { "balance": x.x }
	api.HandleFunc("/wallet/create", h.CreateWallet).Methods(http.MethodPost)

	api.HandleFunc("/wallet/{address}", h.GetWallet).Methods(http.MethodGet)
	api.HandleFunc("/wallet/{address}", h.RemoveWallet).Methods(http.MethodDelete)
	api.HandleFunc("/wallet/{address}/balance", h.UpdateBalance).Methods(http.MethodPut)

	return r
}