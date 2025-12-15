package models

import "time"

type Transaction struct {
	ID          int       `json:"-"`
	AccountFrom int       `json:"account_from,omitempty"`
	AccountTo   int       `json:"account_to,omitempty"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}
type TransferRequest struct {
	ToPhone string  `json:"to_phone"`
	Amount  float64 `json:"amount"`
}
type TransactionHistory struct {
	AccountTo int       `json:"account_to"`
	ToPhone   *string   `json:"to_phone,omitempty"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
