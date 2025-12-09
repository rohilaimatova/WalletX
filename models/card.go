package models

import "time"

type Card struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	AccountID    int       `json:"account_id"`
	MaskedNumber string    `json:"masked_number"` // 1234****0000
	CardType     string    `json:"card_type"`     // например, "VISA" или "MasterCard"
	ExpiryDate   string    `json:"expiry_date"`   // формат MM/YY
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type CreateCardRequest struct {
	UserID     int    `json:"user_id"`
	AccountID  int    `json:"account_id"`
	FullNumber string `json:"full_number"` // полный номер карты
	CardType   string `json:"card_type"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"` // только для проверки
	CardNumber string `json:"card_number"`
}
