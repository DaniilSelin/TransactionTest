package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"net/http"

	"TransactionTest/internal/delivery"
	"TransactionTest/config"
	"TransactionTest/internal/logger"
	"TransactionTest/internal/database"
	"TransactionTest/internal/repository"
	"TransactionTest/internal/service"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	ctx, err := logger.New(ctx)
	if err != nil {
		log.Fatalf("Error create logger: %v", err)
	}

	appLogger := logger.GetLoggerFromCtx(ctx)

	// 1. Загружаем конфиг
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		appLogger.Fatal(ctx, "Error loading config", zap.Error(err))
	}

	// 2. Подключаемся к БД
	dbPool, err := database.InitDB(cfg)
	if err != nil {
		appLogger.Fatal(ctx, "Database connection failed", zap.Error(err))
	}
	defer dbPool.Close()

	// 3. Запускаем миграции
	err = database.RunMigrations(ctx, dbPool)
	if err != nil {
		appLogger.Fatal(ctx, "Migration failed", zap.Error(err))
	}

	// 4. Инициализируем репозитории
	transactionRepo := repository.NewTransactionRepository(dbPool)
	walletRepo := repository.NewWalletRepository(dbPool)

	// 5. Инициализируем сервисы
	transactionService := service.NewTransactionService(transactionRepo, walletRepo)
	walletService := service.NewWalletService(walletRepo)

	// 5.5. Создаем, при необходимости, начальные 10 кошельков
	if flagEmpty, err := walletService.IsEmpty(ctx); err != nil {
		appLogger.Fatal(ctx, "Failed to check if wallets table is empty", zap.Error(err))
	} else if flagEmpty {
		for i := 0; i < 10; i++ {
			address, err := walletService.CreateWallet(ctx, 100)
			if err != nil {
				appLogger.Fatal(ctx, "Failed to create initial wallets", zap.Error(err))
			}
			appLogger.Info(ctx, "Created wallet", zap.String("address", address))
		}
	}

	// 6. Создаём роутер
	router := delivery.NewRouter(ctx, transactionService, walletService)

	// 7. Запускаем сервер
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	appLogger.Info(ctx, "Starting server", zap.String("address", addr))

	srv := &http.Server{
		Addr: addr,
		Handler: router,
	}

	go func () {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			appLogger.Fatal(ctx, "HTTP server ListenAndServe", zap.Error(err))
		}
	}()

	// 8. Завршаем работу сервер (Graceful Shutdown)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	appLogger.Info(ctx, "Shutting down server")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	appLogger.Info(ctx, "Server gracefull stopped")
}
