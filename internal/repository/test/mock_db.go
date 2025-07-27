package test

import (
	"context"
)

type MockDB struct {
	BeginFunc      func(ctx context.Context) (ITx, error)
	ExecFunc       func(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error)
	QueryRowFunc   func(ctx context.Context, sql string, args ...interface{}) Row
	QueryFunc      func(ctx context.Context, sql string, args ...interface{}) (Rows, error)
}

func (m *MockDB) Begin(ctx context.Context) (ITx, error) {
	return m.BeginFunc(ctx)
}
func (m *MockDB) Exec(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error) {
	return m.ExecFunc(ctx, sql, arguments...)
}
func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	return m.QueryRowFunc(ctx, sql, args...)
}
func (m *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	return m.QueryFunc(ctx, sql, args...)
}

type MockTx struct {
	CommitFunc   func(ctx context.Context) error
	RollbackFunc func(ctx context.Context) error
	ExecFunc     func(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error)
	QueryRowFunc func(ctx context.Context, sql string, args ...interface{}) Row
	QueryFunc    func(ctx context.Context, sql string, args ...interface{}) (Rows, error)
}

func (m *MockTx) Commit(ctx context.Context) error                 { return m.CommitFunc(ctx) }
func (m *MockTx) Rollback(ctx context.Context) error               { return m.RollbackFunc(ctx) }
func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error) {
	return m.ExecFunc(ctx, sql, arguments...)
}
func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	return m.QueryRowFunc(ctx, sql, args...)
}
func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	return m.QueryFunc(ctx, sql, args...)
}

type MockRow struct {
	ScanFunc func(dest ...interface{}) error
}
func (m *MockRow) Scan(dest ...interface{}) error { return m.ScanFunc(dest...) }

type MockRows struct {
	NextFunc func() bool
	ScanFunc func(dest ...interface{}) error
	CloseFunc func()
	ErrFunc func() error
}
func (m *MockRows) Next() bool                   { return m.NextFunc() }
func (m *MockRows) Scan(dest ...interface{}) error { return m.ScanFunc(dest...) }
func (m *MockRows) Close()                       { m.CloseFunc() }
func (m *MockRows) Err() error                   { return m.ErrFunc() }

type MockCommandTag struct {
	RowsAffectedFunc func() int64
}
func (m *MockCommandTag) RowsAffected() int64 { return m.RowsAffectedFunc() } 