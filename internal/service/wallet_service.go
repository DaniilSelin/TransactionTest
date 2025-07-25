package service

import (
	"context"
	"fmt"
	"net/http"

	"TransactionTest/internal/domain"
	"TransactionTest/internal/errors"
	"TransactionTest/internal/repository"
	goErrors "errors"
	
	"github.com/google/uuid"
)

type WalletService struct {
	walletRepo IWalletRepository
}

func NewWalletService(walletRepo IWalletRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo}
}

func (ws *WalletService) CreateWallet(ctx context.Context, balance float64) (string, error) {
	address := uuid.New().String()

	if (balance < 0) {
		return "", errors.NewCustomError("Balance cannot be negative", http.StatusBadRequest, nil)
	}

	if err := ws.walletRepo.CreateWallet(ctx, address, balance); err != nil {
		return "", errors.NewCustomError("Failed to create wallet", http.StatusInternalServerError, err)
	}

	return address, nil
}

func (ws *WalletService) IsEmpty(ctx context.Context) (bool, error) {
	isEmpty, err := ws.walletRepo.IsEmpty(ctx)
	if err != nil {
		return false, errors.NewCustomError("Failed to check if wallets table is empty", http.StatusInternalServerError, err)
	}
	return isEmpty, nil
}

func (ws *WalletService) GetBalance(ctx context.Context, address string) (float64, error) {
	balance, err := ws.walletRepo.GetWalletBalance(ctx, address)
	if err != nil {
		if goErrors.Is(err, repository.ErrWalletNotFound) {
			return 0, errors.NewCustomError("Wallet not found", http.StatusNotFound, err)
		}
		return 0, errors.NewCustomError(fmt.Sprintf("Failed to get balance for wallet %s", address), http.StatusInternalServerError, err)
	}
	return balance, nil
}

func (ws *WalletService) GetWallet(ctx context.Context, address string) (*domain.Wallet, error) {
	wallet, err := ws.walletRepo.GetWallet(ctx, address)
	if err != nil {
		if goErrors.Is(err, repository.ErrWalletNotFound) {
			return nil, errors.NewCustomError("Wallet not found", http.StatusNotFound, err)
		}
		return nil, errors.NewCustomError(fmt.Sprintf("Failed to get wallet %s", address), http.StatusInternalServerError, err)
	}

	return wallet, nil
}

func (ws *WalletService) UpdateBalance(ctx context.Context, address string, newBalance float64) error {
	if newBalance < 0 {
		return errors.NewCustomError("Balance cannot be negative", http.StatusBadRequest, nil)
	}
	err := ws.walletRepo.UpdateWalletBalabnce(ctx, address, newBalance)
	if err != nil {
		if goErrors.Is(err, repository.ErrWalletNotFound) {
			return errors.NewCustomError("Wallet not found", http.StatusNotFound, err)
		}
		return errors.NewCustomError("Failed to update wallet balance", http.StatusInternalServerError, err)
	}
	return nil
}

func (ws *WalletService) RemoveWallet(ctx context.Context, address string) error {
	err := ws.walletRepo.RemoveWallet(ctx, address)
	if err != nil {
		if goErrors.Is(err, repository.ErrWalletNotFound) {
			return errors.NewCustomError("Wallet not found", http.StatusNotFound, err)
		}
		return errors.NewCustomError("Failed to remove wallet", http.StatusInternalServerError, err)
	}
	return nil
}

func (ws *WalletService) CreateWalletsForSeeding(
	ctx context.Context,
	count int, 
	balance float64, 
	failOnError bool,
) (<-chan string, <-chan error){ // забыл как вернуть канал
	wallets := make(chan domain.Wallet)
	done := make(chan string)
    	errChan := make(chan error, 1)
	
	go ws.walletRepo.BatchCreateWallets(
		ctx,
		failOnError,
		wallets,
		done,
		errChan,
	)

	go func() {
		defer close(wallets)
		for i := 0; i < count; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				addr := uuid.New().String()
				wallets <- domain.Wallet{Address: addr, Balance: balance}
			}
		}
	}()

	return done, errChan
}