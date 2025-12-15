package models

type ErrorResponse struct {
	Message string      `json:"message" example:"error message"`
	Error   interface{} `json:"error,omitempty"`
}
type SignUpRequest struct {
	Phone string `json:"phone" example:"+992931062345"`
	Code  string `json:"code,omitempty" example:"12345678"` // опционально
}

// RegisterResponse возвращается после успешной регистрации
type RegisterResponse struct {
	Message string `json:"message" example:"user registered"`
	UserID  int    `json:"user_id" example:"6"`
}

// MessageResponse возвращается после отправки кода
type MessageResponse struct {
	Message string `json:"message" example:"registration code sent"`
}
