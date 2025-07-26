package service

import (
	"context"
	"errors"

	"TransactionTest/internal/domain"
	"TransactionTest/internal/logger"
	
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WalletService struct {
	walletRepo IWalletRepository
	log logger.Logger
}

func NewWalletService(wr IWalletRepository, l logger.Logger) *WalletService {
	return &WalletService{
		walletRepo: wr,
		log: l,
	}
}

func (ws *WalletService) CreateWallet(ctx context.Context, balance float64) (string, domain.ErrorCode) {
	if balance < 0 {
		ws.log.Warn(ctx, "CreateWallet: negative balance not allowed")
		return "", domain.CodeNegativeBalance
	}

	address := uuid.New().String()

	if err := ws.walletRepo.CreateWallet(ctx, address, balance); err != nil {
		switch {
		case errors.Is(err, domain.ErrInternal): // Для ускорения проверок
			ws.log.Error(ctx, "CreateWalett",  zap.Error(err))
			return "", domain.CodeInternal // 99% ошибок
		case errors.Is(err, domain.ErrWalletAlreadyExists):
			ws.log.Warn(ctx, "CreateWalett",  zap.Error(err))
			return "", domain.CodeDuplicateWallet
		case errors.Is(err, domain.ErrNegativeBalance): // Никогда не сработает
			ws.log.Warn(ctx, "CreateWalett",  zap.Error(err))
			return "", domain.CodeNegativeBalance
		default:
			ws.log.Error(ctx, "CreateWalett: unexpected",  zap.Error(err))
			return "", domain.CodeInternal
		}
	}
	ws.log.Info(ctx, "CreateWallet: success create wallet", zap.String("address", address))
	return address, domain.CodeOK
}

func (ws *WalletService) GetBalance(ctx context.Context, address string) (float64, domain.ErrorCode) {
	balance, err := ws.walletRepo.GetWalletBalance(ctx, address)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ws.log.Warn(ctx, "GetBalance: wallet not found",  zap.Error(err))
			return 0, domain.CodeWalletNotFound
		}
		ws.log.Error(ctx, "GetBalance",  zap.Error(err))
		return 0, domain.CodeInternal
	}
	ws.log.Info(ctx, "GetBalance: success get wallet", zap.String("address", address), zap.Float64("balance", balance))
	return balance, domain.CodeOK
}

func (ws *WalletService) GetWallet(ctx context.Context, address string) (*domain.Wallet, domain.ErrorCode) {
	wallet, err := ws.walletRepo.GetWallet(ctx, address)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ws.log.Warn(ctx, "GetWallet: wallet not found",  zap.Error(err))
			return nil, domain.CodeWalletNotFound
		}
		ws.log.Error(ctx, "GetBalance",  zap.Error(err))
		return nil, domain.CodeInternal
	}
	ws.log.Info(ctx, "GetWallet: success get wallet",  zap.String("address", address))
	return wallet, domain.CodeOK
}

func (ws *WalletService) UpdateBalance(ctx context.Context, address string, newBalance float64) domain.ErrorCode {
	if newBalance < 0 {
		ws.log.Warn(ctx, "UpdateBalance: negative balance not allowed")
		return domain.CodeNegativeBalance
	}

	err := ws.walletRepo.UpdateWalletBalance(ctx, address, newBalance)
	if err != nil {
			switch {
			case errors.Is(err, domain.ErrNotFound):
				ws.log.Warn(ctx, "UpdateBalance: wallet not found",  zap.Error(err))
				return domain.CodeWalletNotFound
			case errors.Is(err, domain.ErrNegativeBalance): // Никогда не сработает
				ws.log.Warn(ctx, "UpdateBalance",  zap.Error(err))
				return domain.CodeNegativeBalance
			default:
				ws.log.Error(ctx, "UpdateBalance",  zap.Error(err))
				return domain.CodeInternal
		}
	}
	ws.log.Info(ctx, "UpdateBalance: success update wallet", zap.String("address", address), zap.Float64("newBalance", newBalance))
	return domain.CodeOK
}

func (ws *WalletService) RemoveWallet(ctx context.Context, address string) domain.ErrorCode {
	err := ws.walletRepo.RemoveWallet(ctx, address)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ws.log.Warn(ctx, "RemoveWallet: wallet not found",  zap.Error(err))
			return domain.CodeWalletNotFound
		}
		ws.log.Error(ctx, "RemoveWallet",  zap.Error(err))
		return domain.CodeInternal
	}
	ws.log.Info(ctx, "UpdateBalance: success remove wallet", zap.String("address", address))
	return domain.CodeOK
}

func (ws *WalletService) CreateWalletsForSeeding(
	ctx context.Context,
	count int, 
	balance float64, 
	failOnError bool,
) (<-chan string, <-chan error){
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