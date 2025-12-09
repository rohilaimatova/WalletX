package models

import "time"

type Account struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Balance      float64   `json:"balance"`       // основной баланс
	BonusBalance float64   `json:"bonus_balance"` // бонусный баланс
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
