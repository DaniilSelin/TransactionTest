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

func TestTransactionRepository_GetLastTransactions_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.Rows, error) {
			calls := 0
			return &MockRows{
				NextFunc: func() bool {
					calls++
					return calls == 1
				},
				ScanFunc: func(dest ...interface{}) error {
					*dest[0].(*int64) = 1
					*dest[1].(*string) = "from"
					*dest[2].(*string) = "to"
					*dest[3].(*float64) = 10
					*dest[4].(*time.Time) = time.Now()
					return nil
				},
				CloseFunc: func() {},
				ErrFunc:   func() error { return nil },
			}, nil
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	trs, err := repo.GetLastTransactions(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, trs, 1)
	assert.Equal(t, "from", trs[0].From)
}

func TestTransactionRepository_GetLastTransactions_QueryError(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.Rows, error) {
			return nil, errors.New("fail")
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	trs, err := repo.GetLastTransactions(ctx, 1)
	assert.True(t, errors.Is(err, domain.ErrInternal))
	assert.Nil(t, trs)
}

func TestTransactionRepository_GetLastTransactions_ScanError(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.Rows, error) {
			calls := 0
			return &MockRows{
				NextFunc: func() bool {
					calls++
					return calls == 1
				},
				ScanFunc: func(dest ...interface{}) error {
					return errors.New("fail-scan")
				},
				CloseFunc: func() {},
				ErrFunc:   func() error { return nil },
			}, nil
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	trs, err := repo.GetLastTransactions(ctx, 1)
	assert.True(t, errors.Is(err, domain.ErrInternal))
	assert.Nil(t, trs)
}

func TestTransactionRepository_GetLastTransactions_RowsErr(t *testing.T) {
	ctx := context.Background()
	mockDB := &MockDB{
		QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (repository.Rows, error) {
			return &MockRows{
				NextFunc:  func() bool { return false },
				ScanFunc:  func(dest ...interface{}) error { return nil },
				CloseFunc: func() {},
				ErrFunc:   func() error { return errors.New("fail-rows") },
			}, nil
		},
	}
	repo := repository.NewTransactionRepository(mockDB)
	trs, err := repo.GetLastTransactions(ctx, 1)
	assert.True(t, errors.Is(err, domain.ErrInternal))
	assert.Nil(t, trs)
}
