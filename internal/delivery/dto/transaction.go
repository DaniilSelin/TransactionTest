package dto

type SendMoneyRequest struct {
	From   string  `json:"from" validate:"required,uuid4"`
	To     string  `json:"to"    validate:"required,uuid4,nefield=From"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type GetTransactionByInfoRequest struct {
	From      string `json:"from" validate:"required,uuid4"`
	To        string `json:"to"    validate:"required,uuid4,nefield=From"`
	CreatedAt string `json:"created_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

type TransactionResponse struct {
	Id        int64   `json:"id"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"created_at"`
}

type TransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
}
