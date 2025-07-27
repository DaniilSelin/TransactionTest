package postgres

import (
	"context"
	"errors"

	"TransactionTest/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CommandTag
type commandTagAdapter struct {
	tag pgconn.CommandTag
}

func (c commandTagAdapter) RowsAffected() int64 {
	return c.tag.RowsAffected()
}

// DBError
type dbErrorAdapter struct {
	err *pgconn.PgError
}

func (e dbErrorAdapter) Error() string {
	return e.err.Error()
}
func (e dbErrorAdapter) SQLState() string {
	return e.err.SQLState()
}
func (e dbErrorAdapter) ConstraintName() string {
	return e.err.ConstraintName
}

type noRowsError struct{}

func (noRowsError) Error() string            { return "no_rows" }
func (noRowsError) SQLState() string         { return "no_rows" }
func (noRowsError) ConstraintName() string   { return "" }

// Row
type rowAdapter struct {
	row pgx.Row
}

func (r rowAdapter) Scan(dest ...interface{}) error {
    err := r.row.Scan(dest...)
    if err != nil {
        if pgErr, ok := err.(*pgconn.PgError); ok {
            return dbErrorAdapter{pgErr}
        }
        if errors.Is(err, pgx.ErrNoRows) {
            return noRowsError{}
        }
    }
    return err
}

// Rows
type rowsAdapter struct {
	rows pgx.Rows
}

func (r *rowsAdapter) Next() bool {
	return r.rows.Next()
}
func (r *rowsAdapter) Scan(dest ...interface{}) error {
	err := r.rows.Scan(dest...)
	if err != nil {
        // если это Postgres no rows — возвращаем именно эту константу
        if errors.Is(err, pgx.ErrNoRows) {
            return pgx.ErrNoRows
        }
        // если это PgError — оборачиваем
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            return dbErrorAdapter{pgErr}
        }
	}
	return err
}

func (r *rowsAdapter) Close() {
	r.rows.Close()
}
func (r *rowsAdapter) Err() error {
	err := r.rows.Err()
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return dbErrorAdapter{pgErr}
		}
	}
	return err
}

// Tx
type txAdapter struct {
	tx pgx.Tx
}

func (t *txAdapter) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}
func (t *txAdapter) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}
func (t *txAdapter) Exec(ctx context.Context, sql string, arguments ...interface{}) (repository.CommandTag, error) {
	tag, err := t.tx.Exec(ctx, sql, arguments...)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return commandTagAdapter{tag}, dbErrorAdapter{pgErr}
		}
	}
	return commandTagAdapter{tag}, err
}
func (t *txAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) repository.Row {
	return rowAdapter{t.tx.QueryRow(ctx, sql, args...)}
}
func (t *txAdapter) Query(ctx context.Context, sql string, args ...interface{}) (repository.Rows, error) {
	rows, err := t.tx.Query(ctx, sql, args...)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return &rowsAdapter{rows}, dbErrorAdapter{pgErr}
		}
	}
	return &rowsAdapter{rows}, err
}

// IDB
type PoolAdapter struct {
	pool *pgxpool.Pool
}

func NewPoolAdapter(pool *pgxpool.Pool) *PoolAdapter {
	return &PoolAdapter{pool: pool}
}

func (p *PoolAdapter) Begin(ctx context.Context) (repository.ITx, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return &txAdapter{tx}, dbErrorAdapter{pgErr}
		}
	}
	return &txAdapter{tx}, err
}
func (p *PoolAdapter) Exec(ctx context.Context, sql string, arguments ...interface{}) (repository.CommandTag, error) {
	tag, err := p.pool.Exec(ctx, sql, arguments...)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return commandTagAdapter{tag}, dbErrorAdapter{pgErr}
		}
	}
	return commandTagAdapter{tag}, err
}
func (p *PoolAdapter) QueryRow(ctx context.Context, sql string, args ...interface{}) repository.Row {
	return rowAdapter{p.pool.QueryRow(ctx, sql, args...)}
}
func (p *PoolAdapter) Query(ctx context.Context, sql string, args ...interface{}) (repository.Rows, error) {
	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return &rowsAdapter{rows}, dbErrorAdapter{pgErr}
		}
	}
	return &rowsAdapter{rows}, err
}
