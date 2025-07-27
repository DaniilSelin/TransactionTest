package postgres

import (
	"context"
	"fmt"
	"time"

	"TransactionTest/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	// устанавливаем search_path
	connString := cfg.Pool.ConnConfig.ConnString()
	connString = fmt.Sprintf("%s&search_path=%s,public", connString, cfg.Schema)

	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("FATAL: unable to parse pool config: %w", err)
	}

	poolCfg.MaxConns = cfg.Pool.MaxConns
	poolCfg.MinConns = cfg.Pool.MinConns
	poolCfg.MaxConnLifetime = cfg.Pool.MaxConnLifetime
	poolCfg.MaxConnLifetimeJitter = cfg.Pool.MaxConnLifetimeJitter
	poolCfg.MaxConnIdleTime = cfg.Pool.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = cfg.Pool.HealthCheckPeriod

	for attempt := 0; attempt <= cfg.ConnectRetries; attempt++ {
		pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
		if err == nil {
			return pool, nil
		}
		if attempt == cfg.ConnectRetries {
			break
		}
		select {
		case <-time.After(cfg.ConnectRetryDelay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf(
		"FATAL: failed to connect after %d attempts: %w",
		cfg.ConnectRetries+1, err,
	)
}
