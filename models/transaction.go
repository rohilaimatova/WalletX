package models

import "time"

type DepositRequest struct {
	AccountID int     `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type WithdrawRequest struct {
	AccountID int     `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type TransferRequest struct {
	FromID int     `json:"from_id"`
	ToID   int     `json:"to_id"`
	Amount float64 `json:"amount"`
}

type Transaction struct {
	ID          int       `json:"id"`
	AccountFrom int       `json:"account_from,omitempty"`
	AccountTo   int       `json:"account_to,omitempty"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}
