package models

import "time"

type Account struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Balance      float64   `json:"balance"`
	BonusBalance float64   `json:"bonus_balance"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
