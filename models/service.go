package models

import "time"

type Services struct {
	ID          int       `json:"id" example:"1"`
	Name        string    `json:"name" example:"mobile"`
	Description string    `json:"description" example:"mobile services"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PayRequest struct {
	ServiceType string  `json:"service_type" example:"internet"`
	Account     string  `json:"account"`
	Amount      float64 `json:"amount" example:"100"`
}
