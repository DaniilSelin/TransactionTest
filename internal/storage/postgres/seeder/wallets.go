package seeder

import (
	"context"
	"errors"
	"fmt"
	"os"

	"TransactionTest/config"
)

// Ошибки для обработки в вызывающем коде (main)
var (
	ErrDisabled  = errors.New("wallet seeder disabled")
	ErrCompleted = errors.New("seeder already executed")
)

// CreateWalletsForSeeding сигнатруа функции, содержащей бизнес логику "seeding"
type CreateWalletsForSeeding func(context.Context, int, float64, bool) (<-chan string, <-chan error, bool)

// LogError сигнатруа функции, которая будет логировать ошибки при ErrDisabled = false
type LogError func(context.Context, error)

// SeedWallets создает 10 кошельков при первом запуске. Создает маркер файл куда записывает адреса кошельков
func SeedWallets(ctx context.Context, cfg config.WalletsSeedConfig, createStart CreateWalletsForSeeding, log LogError) error {
	if !cfg.Enabled {
		return ErrDisabled
	}

	_, err := os.Stat(cfg.MarkerFile)
	switch {
	case err == nil:
		return ErrCompleted
	case !os.IsNotExist(err):
		return fmt.Errorf("FATAL: marker file check failed: %w", err)
	}

	markerFile, err := os.OpenFile(cfg.MarkerFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("FATAL: create marker failed: %w", err)
	}
	defer markerFile.Close()

	ctx, cancel := context.WithCancel(ctx)
	done, errChan, flag := createStart(
		ctx,
		cfg.Count,
		cfg.Balance,
		cfg.FailOnError,
	)

	if !flag {
		return fmt.Errorf("seeding start failed")
	}

	for {
		select {
		case address, ok := <-done:
			if !ok {
				return nil // Успешное завершение
			}
			if _, err := fmt.Fprintln(markerFile, address); err != nil {
				cancel() // Отправляем сигнал отмены через контекст
				return fmt.Errorf("FATAL: failed to write address: %w", err)
			}

		case err, ok := <-errChan:
			if !ok {
				return nil // errChan закрыт без ошибки — конец
			}
			if cfg.FailOnError {
				return fmt.Errorf("seeding failed: %w", err)
			}
			log(ctx, err)
		}
	}
}
