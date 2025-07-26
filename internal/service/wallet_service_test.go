package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"TransactionTest/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) CreateWallet(ctx context.Context, address string, balance float64) error {
	args := m.Called(ctx, address, balance)
	return args.Error(0)
}

func (m *MockWalletRepository) GetWalletBalance(ctx context.Context, address string) (float64, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockWalletRepository) GetWallet(ctx context.Context, address string) (*domain.Wallet, error) {
	args := m.Called(ctx, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockWalletRepository) UpdateWalletBalance(ctx context.Context, address string, balance float64) error {
	args := m.Called(ctx, address, balance)
	return args.Error(0)
}

func (m *MockWalletRepository) RemoveWallet(ctx context.Context, address string) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockWalletRepository) IsEmpty(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockWalletRepository) BatchCreateWallets(
	ctx context.Context,
	failOnError bool,
	wallets <-chan domain.Wallet,
	done chan<- string,
	errChan chan<- error,
) {
	m.Called(ctx, failOnError, wallets, done, errChan)
}

func TestCreateWalletsForSeeding_Success(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	ws := NewWalletService(mockRepo)

	count := 5
	balance := 100.0

	mockRepo.On("BatchCreateWallets", mock.Anything, false, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		repoWallets := args.Get(2).(<-chan domain.Wallet)
		repoDone := args.Get(3).(chan<- string)
		repoErrChan := args.Get(4).(chan<- error)

		// Симулируем успешную обработку всех кошельков
		go func() {
			defer close(repoDone)
			defer close(repoErrChan)
			for w := range repoWallets {
				select {
				case <-args.Get(0).(context.Context).Done(): // Проверяем отмену контекста
					return
				default:
					repoDone <- w.Address
				}
			}
		}()
	}).Return().Once()

	done, errChan := ws.CreateWalletsForSeeding(context.Background(), count, balance, false)

	var createdAddresses []string
	var receivedErrors []error

	for {
		select {
		case addr, ok := <-done:
			if !ok {
				done = nil
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
		if done == nil && errChan == nil {
			break
		}
	}
EndLoop:

	assert.Len(t, createdAddresses, count, "Ожидалось создание %d кошельков", count)
	assert.Empty(t, receivedErrors, "Ожидались отсутствие ошибок")
	mockRepo.AssertExpectations(t)
}

func TestCreateWalletsForSeeding_FailOnErrorTrue(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	ws := NewWalletService(mockRepo)

	count := 5
	balance := 100.0

	mockRepo.On("BatchCreateWallets", mock.Anything, true, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		repoWallets := args.Get(2).(<-chan domain.Wallet)
		repoDone := args.Get(3).(chan<- string)
		repoErrChan := args.Get(4).(chan<- error)

		go func() {
			defer close(repoDone)
			defer close(repoErrChan)
			// Обрабатываем один кошелек успешно, затем симулируем ошибку и останавливаемся
			w := <-repoWallets
			select {
			case <-args.Get(0).(context.Context).Done():
				return
			default:
				repoDone <- w.Address // Первый кошелек успешно
			}

			repoErrChan <- fmt.Errorf("симулированная ошибка базы данных для failOnError=true")
		}()
	}).Return().Once()

	done, errChan := ws.CreateWalletsForSeeding(context.Background(), count, balance, true)

	var createdAddresses []string
	var receivedErrors []error

	for {
		select {
		case addr, ok := <-done:
			if !ok {
				done = nil
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
		if done == nil && errChan == nil {
			break
		}
	}
EndLoop:

	assert.Len(t, createdAddresses, 1, "Ожидалось, что будет создан только 1 кошелек до ошибки")
	assert.Len(t, receivedErrors, 1, "Ожидалась одна ошибка")
	assert.Contains(t, receivedErrors[0].Error(), "симулированная ошибка базы данных для failOnError=true", "Ожидалось конкретное сообщение об ошибке")
	mockRepo.AssertExpectations(t)
}

func TestCreateWalletsForSeeding_FailOnErrorFalse(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	ws := NewWalletService(mockRepo)

	count := 5
	balance := 100.0

	mockRepo.On("BatchCreateWallets", mock.Anything, false, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		repoWallets := args.Get(2).(<-chan domain.Wallet)
		repoDone := args.Get(3).(chan<- string)
		repoErrChan := args.Get(4).(chan<- error)

		go func() {
			defer close(repoDone)
			defer close(repoErrChan)
			processedCount := 0
			for w := range repoWallets {
				select {
				case <-args.Get(0).(context.Context).Done():
					return
				default:
					processedCount++
					if processedCount == 3 { // Симулируем сбой для 3-го кошелька
						// В реальном BatchCreateWallets он бы записал ошибку и продолжил, не отправляя в errChan
						// Поэтому для этого мока мы просто пропускаем отправку в done для этого кошелька
						fmt.Printf("Мок: Пропускаем кошелек %s из-за ошибки (failOnError=false)\n", w.Address)
						continue
					}
					repoDone <- w.Address
				}
			}
		}()
	}).Return().Once()

	done, errChan := ws.CreateWalletsForSeeding(context.Background(), count, balance, false)

	var createdAddresses []string
	var receivedErrors []error

	for {
		select {
		case addr, ok := <-done:
			if !ok {
				done = nil
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
		if done == nil && errChan == nil {
			break
		}
	}
EndLoop:

	assert.Len(t, createdAddresses, count-1, "Ожидалось, что будет создано %d кошельков (один пропущен)", count-1)
	assert.Empty(t, receivedErrors, "Ожидались отсутствие фатальных ошибок")
	mockRepo.AssertExpectations(t)
}

func TestCreateWalletsForSeeding_ContextCancellation(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	ws := NewWalletService(mockRepo)

	count := 10
	balance := 100.0

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	mockRepo.On("BatchCreateWallets", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		repoWallets := args.Get(2).(<-chan domain.Wallet)
		repoDone := args.Get(3).(chan<- string)
		repoErrChan := args.Get(4).(chan<- error)
		repoCtx := args.Get(0).(context.Context)

		go func() {
			processedCount := 0
			for w := range repoWallets {
				select {
				case <-repoCtx.Done():
					fmt.Println("Мок BatchCreateWallets: Контекст отменен, остановка обработки.")
					repoErrChan <- repoCtx.Err() // Симулируем ошибку отката для перехвата
					close(repoDone)
					close(repoErrChan)
					return
				default:
					repoDone <- w.Address
					processedCount++
					time.Sleep(10 * time.Millisecond)
				}
			}
			close(repoDone)
			close(repoErrChan)
		}()
	}).Return().Once()

	done, errChan := ws.CreateWalletsForSeeding(ctx, count, balance, false)

	var createdAddresses []string
	var receivedErrors []error

	for {
		select {
		case addr, ok := <-done:
			if !ok {
				done = nil
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
		if done == nil && errChan == nil {
			break
		}
	}
EndLoop:
	// Поскольку отмена происходит в середине процесса, мы ожидаем меньшее количество кошельков, чем `count`
	assert.Less(t, len(createdAddresses), count, "Ожидалось меньшее количество кошельков (%d) из-за отмены", count)
	assert.Greater(t, len(createdAddresses), 0, "Ожидалось, что будет создан хотя бы один кошелек до отмены")
	assert.Len(t, receivedErrors, 1, "Ожидалась одна ошибка из-за отмены контекста")
	assert.Equal(t, context.Canceled, receivedErrors[0], "Ожидалась ошибка context.Canceled")
	mockRepo.AssertExpectations(t)
}