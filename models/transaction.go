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
	ToPhone string  `json:"to_phone" example:"+992931753756"`
	Amount  float64 `json:"amount" example:"100.34"`
}
type TransactionHistory struct {
	AccountTo int       `json:"account_to" example:"3"`
	ToPhone   *string   `json:"to_phone,omitempty" example:"+992931753756"`
	Amount    float64   `json:"amount" example:"100"`
	Type      string    `json:"type" example:"transfer"`
	CreatedAt time.Time `json:"created_at"`
}
