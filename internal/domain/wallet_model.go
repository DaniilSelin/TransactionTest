package domain

import (
	"time"
)

type Wallet struct {
	Address   string
	Balance   float64
	CreatedAt time.Time
}
