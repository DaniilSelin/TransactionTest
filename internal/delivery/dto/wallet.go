package dto

type CreateWalletRequest struct {
    Balance float64 `json:"balance" validate:"required,gte=0"`
}

type UpdateBalanceRequest struct {
    Balance float64 `json:"balance" validate:"required,gte=0"`
}

type UpdateBalanceRequest struct {
    Balance float64 `json:"balance" validate:"required,gte=0"`
}

type CreateWalletResponse struct {
    Address string `json:"address"`
}

type BalanceResponse struct {
    Balance float64 `json:"balance"`
}

type WalletResponse struct {
    Address   string  `json:"address"`
    Balance   float64 `json:"balance"`
    CreatedAt string  `json:"created_at"`
}