package main

import (
	"context"
	"fmt"
	"log"

	"TransactionTest/config"
	"TransactionTest/migrations"
	"TransactionTest/internal/logger"
	"TransactionTest/internal/storage/postgres"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func main() {
	// Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Выводим содержимое конфига
	fmt.Println("--- Application Configuration ---")
	fmt.Printf("%+v\n", cfg)
	fmt.Println("--- PostgreSQL ConnConfig Details ---")
	fmt.Printf("%+v\n", cfg.Postgres.Pool.ConnConfig)
	fmt.Println("-------------------------------------")

	// Создаем логгер
	appLogger, err := logger.New(&cfg.Logger.Logger)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	
	// Создаем контекст с RequestID
	ctx := context.WithValue(context.Background(), logger.RequestID, uuid.New().String())

	// Тестируем методы логгера
	appLogger.Info(ctx, "This is an info message with RequestID", zap.String("component", "check"))
	appLogger.Warn(ctx, "This is a warning message with RequestID", zap.Int("attempt", 1))
	appLogger.Error(ctx, "This is an error message with RequestID", zap.Error(fmt.Errorf("something went wrong")))

	log.Println("Logger tests completed.")
	// Создаем пулл
	pool, err := postgres.Connect(ctx, &cfg.Postgres)
	if err != nil {
		log.Fatalf("Error creating pool: %v", err)
	}
	fmt.Printf("%+v\n", pool)

	log.Println("ПОПЫТКА МИГРАЦИЙ")
	// Запускаем миграции
	connStr := cfg.Postgres.Pool.ConnConfig.ConnString()
	sourceURL := fmt.Sprintf(
		"%s://%s?schema=%s", cfg.Migrations.Driver ,cfg.Migrations.Dir, cfg.Postgres.Schema,
	)

	fmt.Println(connStr,"-111-", sourceURL)
	m, err := migrations.New(sourceURL, connStr)
	if err != nil {
		log.Fatalf("FATAL: failed to create migrate instance: %v", err)
	}

	appLogger.Info(ctx, "Starting migrations...")
	if err := m.Up(); err != nil {
		if err == migrations.ErrNoChange {
			appLogger.Info(ctx, "No new migrations to apply.")
		} else {
			appLogger.Fatal(ctx, fmt.Sprintf("FATAL: failed to apply migrations: %v", err))
		}
	} else {
		appLogger.Info(ctx, "Migrations applied successfully.")
	}
}
