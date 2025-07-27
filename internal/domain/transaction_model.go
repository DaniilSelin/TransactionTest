package domain

import (
	"time"
)

type Transaction struct {
	Id       int64
	From      string
	To        string
	Amount    float64
	CreatedAt time.Time
}