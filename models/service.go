package models

import "time"

type Services struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PayRequest struct {
	ServiceType string  `json:"service_type"`
	Account     string  `json:"account"`
	Amount      float64 `json:"amount"`
}
