package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"TransactionTest/config"
	httpCust "TransactionTest/internal/delivery/http"
	"TransactionTest/internal/delivery/http/handler"
	"TransactionTest/internal/logger"
	"TransactionTest/internal/repository"
	"TransactionTest/internal/service"
	"TransactionTest/internal/storage/postgres"
	"TransactionTest/internal/storage/postgres/seeder"
	"TransactionTest/migrations"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
	config     config.Config
}

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger, err := logger.New(&cfg.Logger.Logger)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	appLogger.Info(ctx, "Starting application...")

	server, err := NewServer(cfg, appLogger)
	if err != nil {
		appLogger.Fatal(ctx, fmt.Sprintf("Failed to create server: %v", err))
	}

	// Запускаем сервер
	go func() {
		appLogger.Info(ctx, fmt.Sprintf("Server starting on %s:%d", cfg.Server.Host, cfg.Server.Port))
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal(ctx, fmt.Sprintf("Server failed to start: %v", err))
		}
	}()

	server.WaitForShutdown(ctx)
}

// NewServer создает новый экземпляр сервера
func NewServer(cfg config.Config, appLogger logger.Logger) (*Server, error) {
	ctx := context.WithValue(context.Background(), logger.RequestID, uuid.New().String())

	pool, err := postgres.Connect(ctx, &cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	appLogger.Info(ctx, "Database connection established")

	if err := runMigrations(ctx, cfg, appLogger); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	adapter := postgres.NewPoolAdapter(pool)

	transactionRepo := repository.NewTransactionRepository(adapter)
	walletRepo := repository.NewWalletRepository(adapter)

	transactionService := service.NewTransactionService(transactionRepo, walletRepo, appLogger)
	walletService := service.NewWalletService(walletRepo, appLogger)

	if err = Seeding(ctx, cfg.Seeding.Wallets, walletService.CreateWalletsForSeeding, appLogger); err != nil {
		return nil, fmt.Errorf("Seeding failed: %w", err)
	}

	h := handler.NewHandler(transactionService, walletService, appLogger)

	r := httpCust.NewRouter(h, appLogger)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		logger:     appLogger,
		config:     cfg,
	}, nil
}

// Start запускает HTTP сервер
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// WaitForShutdown ожидает сигнал для graceful shutdown
func (s *Server) WaitForShutdown(ctx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ждем сигнал
	<-quit
	s.logger.Info(ctx, "Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Graceful shutdown
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logger.Error(ctx, "Server forced to shutdown", zap.Error(err))
	} else {
		s.logger.Info(ctx, "Server gracefully stopped")
	}
}

// Seeding запускает процесс создания стартового количесвта кошельков
func Seeding(ctx context.Context, cfg config.WalletsSeedConfig, createStart seeder.CreateWalletsForSeeding, appLogger logger.Logger) error {
	seedlog := func(ctx context.Context, err error) {
		appLogger.Warn(ctx, "Seeding error create wallet", zap.Error(err))
	}

	appLogger.Info(ctx, "Start seeding...")

	err := seeder.SeedWallets(ctx, cfg, createStart, seedlog)
	if err != nil {
		switch err {
		case seeder.ErrDisabled:
			appLogger.Warn(ctx, "seeding disable")
			return nil
		case seeder.ErrCompleted:
			appLogger.Warn(ctx, "seeder already executed")
			return nil
		default:
			appLogger.Fatal(ctx, fmt.Sprintf("FATAL: failed to run seed: %v", err))
			return err
		}
	}
	return nil
}

// runMigrations запускает миграции базы данных
func runMigrations(ctx context.Context, cfg config.Config, appLogger logger.Logger) error {
	// Формируем строки подключения
	connStr := cfg.Postgres.Pool.ConnConfig.ConnString()
	sourceURL := fmt.Sprintf(
		"%s://%s?schema=%s",
		cfg.Migrations.Driver,
		cfg.Migrations.Dir,
		cfg.Postgres.Schema,
	)

	appLogger.Info(ctx, "Initializing migrations...",
		zap.String("source", sourceURL),
		zap.String("database", cfg.Postgres.Pool.ConnConfig.Database))

	m, err := migrations.New(sourceURL, connStr)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	appLogger.Info(ctx, "Starting migrations...")
	if err := m.Up(); err != nil {
		if err == migrations.ErrNoChange {
			appLogger.Info(ctx, "No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	appLogger.Info(ctx, "Migrations applied successfully")
	return nil
}
