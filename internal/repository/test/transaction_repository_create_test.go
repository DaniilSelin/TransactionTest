package test

import (
	"TransactionTest/internal/domain"
	"TransactionTest/internal/repository"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockDBError struct {
	sqlState   string
	constraint string
}

func (e *mockDBError) Error() string          { return "mock db error" }
func (e *mockDBError) SQLState() string       { return e.sqlState }
func (e *mockDBError) ConstraintName() string { return e.constraint }

func TestTransactionRepository_CreateTransaction_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				*dest[0].(*int64) = 42
				return nil
			}}
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	id, err := repo.CreateTransaction(ctx, "from", "to", 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), id)
}

func TestTransactionRepository_CreateTransaction_NegativeAmount(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return &mockDBError{sqlState: repository.ErrCodeCheckViolation, constraint: repository.ConstraintAmountPositive}
			}}
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	_, err := repo.CreateTransaction(ctx, "from", "to", -10)
	assert.True(t, errors.Is(err, domain.ErrNegativeAmount))
}

func TestTransactionRepository_CreateTransaction_SelfTransfer(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return &mockDBError{sqlState: repository.ErrCodeCheckViolation, constraint: repository.ConstraintNoSelfTransfer}
			}}
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	_, err := repo.CreateTransaction(ctx, "from", "from", 10)
	assert.True(t, errors.Is(err, domain.ErrSelfTransfer))
}

func TestTransactionRepository_CreateTransaction_FKViolation(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return &mockDBError{sqlState: repository.ErrCodeForeignKeyViolation}
			}}
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	_, err := repo.CreateTransaction(ctx, "from", "to", 10)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestTransactionRepository_CreateTransaction_InternalError(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return errors.New("fail")
			}}
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	_, err := repo.CreateTransaction(ctx, "from", "to", 10)
	assert.True(t, errors.Is(err, domain.ErrInternal))
}

func TestTransactionRepository_CreateTransactionTx_Success(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				*dest[0].(*int64) = 42
				return nil
			}}
		},
	}
	repo := repository.NewTransactionRepository(nil)
	id, err := repo.CreateTransactionTx(ctx, mockTx, "from", "to", 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), id)
}

func TestTransactionRepository_CreateTransactionTx_NegativeAmount(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return &mockDBError{sqlState: repository.ErrCodeCheckViolation, constraint: repository.ConstraintAmountPositive}
			}}
		},
	}
	repo := repository.NewTransactionRepository(nil)
	_, err := repo.CreateTransactionTx(ctx, *mockTx, "from", "to", -10)
	assert.True(t, errors.Is(err, domain.ErrNegativeAmount))
}

func TestTransactionRepository_CreateTransactionTx_SelfTransfer(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return &mockDBError{sqlState: repository.ErrCodeCheckViolation, constraint: repository.ConstraintNoSelfTransfer}
			}}
		},
	}
	repo := repository.NewTransactionRepository(nil)
	_, err := repo.CreateTransactionTx(ctx, *mockTx, "from", "from", 10)
	assert.True(t, errors.Is(err, domain.ErrSelfTransfer))
}

func TestTransactionRepository_CreateTransactionTx_FKViolation(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return &mockDBError{sqlState: repository.ErrCodeForeignKeyViolation}
			}}
		},
	}
	repo := repository.NewTransactionRepository(nil)
	_, err := repo.CreateTransactionTx(ctx, *mockTx, "from", "to", 10)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestTransactionRepository_CreateTransactionTx_InternalError(t *testing.T) {
	ctx := context.Background()
	mockTx := &MockTx{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return errors.New("fail")
			}}
		},
	}
	repo := repository.NewTransactionRepository(nil)
	_, err := repo.CreateTransactionTx(ctx, *mockTx, "from", "to", 10)
	assert.True(t, errors.Is(err, domain.ErrInternal))
}
