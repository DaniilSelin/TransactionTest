package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"TransactionTest/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
func (m *MockTx) Rollback(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error) {
    callArgs := m.Called(ctx, sql, arguments)
    return callArgs.Get(0).(mockCommandTag), callArgs.Error(1)
}
func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
    callArgs := m.Called(ctx, sql, args)
    return callArgs.Get(0).(Rows), callArgs.Error(1)
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
    callArgs := m.Called(ctx, sql, args)
    return callArgs.Get(0).(Row)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin(ctx context.Context) (Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(Tx), args.Error(1)
}
func (m *MockDB) Exec(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	return args.Get(0).(mockCommandTag), args.Error(1)
}
func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	return m.Called(ctx, sql, args).Get(0).(Row)
}
func (m *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	mockArgs := m.Called(ctx, sql, args)
    return mockArgs.Get(0).(Rows), mockArgs.Error(1)
}

type mockCommandTag struct {
    rowsAffected int64
}

func (m mockCommandTag) RowsAffected() int64 {
    return m.rowsAffected
}

func TestBatchCreateWallets_Success(t *testing.T) {
	mockDB := new(MockDB)
	mockTx := new(MockTx)
	walletRepo := NewWalletRepository(mockDB)

	count := 5
	balance := 100.0

	walletsChan := make(chan domain.Wallet)
	doneChan := make(chan string)
	errChan := make(chan error, 1)

	mockDB.On("Begin", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockCommandTag{1}, nil).Times(count)
	mockTx.On("Commit", mock.Anything).Return(nil).Once()

	go walletRepo.BatchCreateWallets(context.Background(), false, walletsChan, doneChan, errChan)

	go func() {
		defer close(walletsChan)
		for i := 0; i < count; i++ {
			walletsChan <- domain.Wallet{Address: uuid.New().String(), Balance: balance}
		}
	}()

	var createdAddresses []string
	var receivedErrors []error

	for {
		select {
		case addr, ok := <-doneChan:
			if !ok {
				doneChan = nil
				continue
			}
			createdAddresses = append(createdAddresses, addr)
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			receivedErrors = append(receivedErrors, err)
		case <-time.After(500 * time.Millisecond):
			goto EndLoop
		}
		if doneChan == nil && errChan == nil {
			break
		}
	}
EndLoop:

	assert.Len(t, createdAddresses, count, "Ожидалось создание %d кошельков", count)
	assert.Empty(t, receivedErrors, "Ожидались отсутствие ошибок")
	mockDB.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestBatchCreateWallets_FailOnErrorTrue(t *testing.T) {
	mockDB := new(MockDB)
	mockTx := new(MockTx)
	walletRepo := NewWalletRepository(mockDB)

	count := 5
	balance := 100.0

	walletsChan := make(chan domain.Wallet)
	doneChan := make(chan string)
	errChan := make(chan error, 1)

	mockDB.On("Begin", mock.Anything).Return(mockTx, nil).Once()
	// Успешная вставка для первых 2 кошельков
	mockTx.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockCommandTag{1}, nil).Once()
	mockTx.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockCommandTag{1}, nil).Once()
	// Затем ошибка
	mockTx.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockCommandTag{0}, errors.New("mocked db error")).Once()
	mockTx.On("Rollback", mock.Anything).Return(nil).Once()

	go walletRepo.BatchCreateWallets(context.Background(), true, walletsChan, doneChan, errChan)

	go func() {
		defer close(walletsChan)
		for i := 0; i < count; i++ {
			walletsChan <- domain.Wallet{Address: uuid.New().String(), Balance: balance}
		}
	}()

	var createdAddresses []string
	var receivedErrors []error

	for {
		select {
		case addr, ok := <-doneChan:
			if !ok {
				doneChan = nil
				continue
			}
			createdAddresses = append(createdAddresses, addr)
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			receivedErrors = append(receivedErrors, err)
		case <-time.After(500 * time.Millisecond):
			goto EndLoop
		}
		if doneChan == nil && errChan == nil {
			break
		}
	}
EndLoop:

	assert.Len(t, createdAddresses, 2, "Expected 2 wallets to be created before error")
	assert.Len(t, receivedErrors, 1, "Expected one error")
	assert.Contains(t, receivedErrors[0].Error(), "mocked db error", "Expected specific error message")
	mockDB.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestBatchCreateWallets_FailOnErrorFalse(t *testing.T) {
	mockDB := new(MockDB)
	mockTx := new(MockTx)
	walletRepo := NewWalletRepository(mockDB)

	count := 5
	balance := 100.0

	walletsChan := make(chan domain.Wallet)
	doneChan := make(chan string)
	errChan := make(chan error, 1)

	mockDB.On("Begin", mock.Anything).Return(mockTx, nil).Once()
	// Успешная вставка для первых 2 кошельков
	mockTx.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockCommandTag{1}, nil).Once()
	mockTx.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockCommandTag{1}, nil).Once()
	// Затем ошибка
	mockTx.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockCommandTag{0}, errors.New("mocked db error")).Once()
	// Затем снова успешные
	mockTx.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockCommandTag{1}, nil).Once()
	mockTx.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(mockCommandTag{1}, nil).Once()
	mockTx.On("Commit", mock.Anything).Return(nil).Once()

	go walletRepo.BatchCreateWallets(context.Background(), false, walletsChan, doneChan, errChan)

	go func() {
		defer close(walletsChan)
		for i := 0; i < count; i++ {
			walletsChan <- domain.Wallet{Address: uuid.New().String(), Balance: balance}
		}
	}()

	var createdAddresses []string
	var receivedErrors []error

	for {
		select {
		case addr, ok := <-doneChan:
			if !ok {
				doneChan = nil
				continue
			}
			createdAddresses = append(createdAddresses, addr)
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			receivedErrors = append(receivedErrors, err)
		case <-time.After(500 * time.Millisecond):
			goto EndLoop
		}
		if doneChan == nil && errChan == nil {
			break
		}
	}
EndLoop:

	assert.Len(t, createdAddresses, count-1, "Expected %d wallets to be created (one skipped)", count-1)
	mockDB.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestBatchCreateWallets_ContextCancellation(t *testing.T) {
	mockDB := new(MockDB)
	mockTx := new(MockTx)
	walletRepo := NewWalletRepository(mockDB)

	count := 5
	balance := 100.0

	ctx, cancel := context.WithCancel(context.Background())

	walletsChan := make(chan domain.Wallet)
	doneChan := make(chan string)
	errChan := make(chan error, 1)

	mockDB.On("Begin", mock.Anything).Return(mockTx, nil).Once()
	// Мокируем Exec так, чтобы он зависал, пока контекст не будет отменен
	mockTx.On("Exec", mock.Anything, mock.Anything, mock.Anything).Run(func(mockArgs mock.Arguments) {
		// Блокируем, пока контекст не будет отменен
		<-mockArgs.Get(0).(context.Context).Done()
	}).Return(mockCommandTag{0}, context.Canceled).Maybe()
	
	mockTx.On("Rollback", mock.Anything).Return(nil).Once()

	go walletRepo.BatchCreateWallets(ctx, true, walletsChan, doneChan, errChan)

	go func() {
		defer close(walletsChan)
		for i := 0; i < count; i++ {
			walletsChan <- domain.Wallet{Address: uuid.New().String(), Balance: balance}
			time.Sleep(10 * time.Millisecond) // Даем время для обработки
		}
	}()

	// Отменяем контекст через некоторое время
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	var createdAddresses []string
	var receivedErrors []error

	for {
		select {
		case addr, ok := <-doneChan:
			if !ok {
				doneChan = nil
				continue
			}
			createdAddresses = append(createdAddresses, addr)
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			receivedErrors = append(receivedErrors, err)
		case <-time.After(2 * time.Second):
			goto EndLoop
		}
		if doneChan == nil && errChan == nil {
			break
		}
	}
EndLoop:

	assert.GreaterOrEqual(t, len(createdAddresses), 0, "Expected some wallets to be created or none if cancelled early")
	assert.Less(t, len(createdAddresses), count, "Expected fewer than %d wallets due to cancellation", count)
	assert.Len(t, receivedErrors, 1, "Expected one error due to context cancellation")
	assert.Contains(t, receivedErrors[0].Error(), "context canceled", "Expected context.Canceled error")
	mockDB.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	
}

func TestBatchCreateWallets_ContextAlreadyCancelled(t *testing.T) {
	mockDB := new(MockDB)
	mockTx := new(MockTx)
	walletRepo := NewWalletRepository(mockDB)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст сразу же

	walletsChan := make(chan domain.Wallet)
	doneChan := make(chan string)
	errChan := make(chan error, 1)

	mockDB.On("Begin", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Rollback", mock.Anything).Return(nil).Maybe() // Может не вызваться, если цикл не начался
	// Мокируем Commit, чтобы он возвращал ошибку отмены контекста, так как транзакция будет коммититься на отмененном контексте
	mockTx.On("Commit", mock.Anything).Return(context.Canceled).Once()

	go walletRepo.BatchCreateWallets(ctx, true, walletsChan, doneChan, errChan)

	// Закрываем канал walletsChan немедленно, чтобы цикл в BatchCreateWallets завершился
	close(walletsChan)

	var createdAddresses []string
	var receivedErrors []error

	for {
		select {
		case addr, ok := <-doneChan:
			if !ok {
				doneChan = nil
				continue
			}
			createdAddresses = append(createdAddresses, addr)
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			receivedErrors = append(receivedErrors, err)
		case <-time.After(500 * time.Millisecond):
			goto EndLoop
		}
		if doneChan == nil && errChan == nil {
			break
		}
	}
EndLoop:

	assert.Empty(t, createdAddresses, "Ожидалось отсутствие созданных кошельков")
	assert.Len(t, receivedErrors, 1, "Ожидалась одна ошибка")
	assert.Contains(t, receivedErrors[0].Error(), "context canceled", "Ожидалась ошибка отмены контекста")
	mockDB.AssertExpectations(t)
	mockTx.AssertExpectations(t)
} 