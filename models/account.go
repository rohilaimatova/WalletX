package models

import "time"

type CreateAccountRequest struct {
	UserID int `json:"user_id"`
}

type Account struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
type BalanceResponse struct {
	AccountID int     `json:"account_id"`
	Balance   float64 `json:"balance"`
}
