package postgres

import (
	"reflect"
	"unsafe"
	"context"
	"time"
	"fmt"

	"TransactionTest/config"

	"github.com/jackc/pgx/v5/pgxpool"
)
/*
func forceSetConnString(cfg *pgxpool.Config, value string) {
	v := reflect.ValueOf(cfg).Elem()
	f := v.FieldByName("connString")
	if !f.IsValid() {
		panic("connString field not found in pgxpool.Config")
	}
	ptr := unsafe.Pointer(f.UnsafeAddr())
	reflect.NewAt(f.Type(), ptr).Elem().SetString(value)
}*/

// createdByParseConfig = true
func forceSetCreatedByParseConfig(cfg *pgxpool.Config) {
	v := reflect.ValueOf(cfg).Elem()
	f := v.FieldByName("createdByParseConfig")
	if !f.IsValid() {
		panic("createdByParseConfig field not found in pgxpool.Config")
	}
	ptr := unsafe.Pointer(f.UnsafeAddr())
	reflect.NewAt(f.Type(), ptr).Elem().SetBool(true)
}

func Connect(ctx context.Context, poolCfg *pgxpool.Config, migCfg *config.MigrationConfig) (*pgxpool.Pool, error) {
    // dsn := poolCfg.ConnConfig.ConnString()

	forceSetCreatedByParseConfig(poolCfg)

    var err error
    for attempt := 0; attempt <= migCfg.ConnectRetries; attempt++ {
        pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
        if err == nil {
            return pool, nil
        }
        if attempt == migCfg.ConnectRetries {
            break
        }
        select {
        case <-time.After(migCfg.ConnectRetryDelay):
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }

    return nil, fmt.Errorf(
        "FATAL: failed to connect after %d attempts: %w",
        migCfg.ConnectRetries+1, err,
    )
}