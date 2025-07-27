package repository

import (
	"context"
)

type IDB interface {
	Begin(ctx context.Context) (ITx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)
}

type ITx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)
}

type Row interface {
	Scan(dest ...interface{}) error
}

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close()
	Err() error
}

type CommandTag interface {
	RowsAffected() int64
}

// Абстракция ошибки, возвращаемой драйвером БД
// Позволяет обрабатывать ошибки целостности и ограничения
type DBError interface {
	error
	SQLState() string
	ConstraintName() string
}

// Коды ошибок для проверки в репозитории (Postgres SQLSTATE)
const (
	ErrCodeUniqueViolation     = "23505"
	ErrCodeCheckViolation      = "23514"
	ErrCodeForeignKeyViolation = "23503"
)

// Имена constraint-ов
const (
	ConstraintBalanceNonNegative = "chk_balance_nonnegative"
	ConstraintAmountPositive     = "chk_amount_positive"
	ConstraintNoSelfTransfer     = "chk_no_self_transfer"
)
