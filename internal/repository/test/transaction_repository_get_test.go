package test

import (
	"TransactionTest/internal/domain"
	"TransactionTest/internal/repository"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransactionRepository_GetTransactionById_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				*dest[0].(*int64) = 1
				*dest[1].(*string) = "from"
				*dest[2].(*string) = "to"
				*dest[3].(*float64) = 10
				*dest[4].(*time.Time) = time.Now()
				return nil
			}}
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	tr, err := repo.GetTransactionById(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, tr)
	assert.Equal(t, int64(1), tr.Id)
}

func TestTransactionRepository_GetTransactionById_NotFound(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return &mockDBError{sqlState: "no_rows"}
			}}
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	tr, err := repo.GetTransactionById(ctx, 1)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	assert.Nil(t, tr)
}

func TestTransactionRepository_GetTransactionById_InternalError(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return errors.New("fail")
			}}
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	tr, err := repo.GetTransactionById(ctx, 1)
	assert.True(t, errors.Is(err, domain.ErrInternal))
	assert.Nil(t, tr)
}

func TestTransactionRepository_GetTransactionByInfo_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				*dest[0].(*int64) = 1
				*dest[1].(*string) = "from"
				*dest[2].(*string) = "to"
				*dest[3].(*float64) = 10
				*dest[4].(*time.Time) = time.Now()
				return nil
			}}
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	tr, err := repo.GetTransactionByInfo(ctx, "from", "to", time.Now())
	assert.NoError(t, err)
	assert.NotNil(t, tr)
	assert.Equal(t, "from", tr.From)
}

func TestTransactionRepository_GetTransactionByInfo_NotFound(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return &mockDBError{sqlState: "no_rows"}
			}}
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	tr, err := repo.GetTransactionByInfo(ctx, "from", "to", time.Now())
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	assert.Nil(t, tr)
}

func TestTransactionRepository_GetTransactionByInfo_InternalError(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) repository.Row {
			return &MockRow{ScanFunc: func(dest ...interface{}) error {
				return errors.New("fail")
			}}
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	tr, err := repo.GetTransactionByInfo(ctx, "from", "to", time.Now())
	assert.True(t, errors.Is(err, domain.ErrInternal))
	assert.Nil(t, tr)
}
