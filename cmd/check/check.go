package main

import (
	"context"
	"fmt"
	"log"

	"TransactionTest/config"
	"TransactionTest/internal/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func main() {
	// 1. Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 2. Выводим содержимое конфига
	fmt.Println("--- Application Configuration ---")
	fmt.Printf("%+v\n", cfg)
	fmt.Println("---------------------------------")

	// 3. Создаем логгер
	appLogger, err := logger.New(&cfg.LoggerConfig.Logger)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	// Создаем контекст с RequestID
	ctx := context.WithValue(context.Background(), logger.RequestID, uuid.New().String())

	// 4. Тестируем методы логгера
	appLogger.Info(ctx, "This is an info message with RequestID", zap.String("component", "check"))
	appLogger.Warn(ctx, "This is a warning message with RequestID", zap.Int("attempt", 1))
	appLogger.Error(ctx, "This is an error message with RequestID", zap.Error(fmt.Errorf("something went wrong")))

	log.Println("Logger tests completed.")
}
