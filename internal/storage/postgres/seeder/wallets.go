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
	ErrDisabled  = errors.New("WARN: wallet seeder disabled")
	ErrCompleted = errors.New("WARN: seeder already executed")
)

type CreateWalletsForSeeding func(context.Context, int, float64, bool) (<-chan string, <-chan error) 

func SeedWallets(ctx context.Context, cfg config.WalletsSeedConfig, createStart CreateWalletsForSeeding) error {
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
	done, errChan := createStart(
		ctx,
		cfg.Count,
		cfg.Balance,
		cfg.FailOnError,
	)

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
			return fmt.Errorf("seeding failed: %w", err)
		}
	}
}