package models

type User struct {
	ID               int     `json:"id"`
	Phone            string  `json:"phone"`
	Password         string  `json:"-"` // не отдаём наружу
	FirstName        *string `json:"first_name,omitempty"`
	LastName         *string `json:"last_name,omitempty"`
	MiddleName       *string `json:"middle_name,omitempty"`
	PassportNumber   *string `json:"passport_number,omitempty"`
	IsVerified       bool    `json:"is_verified"`
	DeviceID         bool    `json:"device_id,omitempty"`
	PasswordAttempts int     `json:"password_attempts" db:"password_attempts"`
}
type RegisterRequest struct {
	Phone string `json:"phone"`
}
type SetPasswordRequest struct {
	UserID   int    `json:"user_id"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
type VerifyIdentityRequest struct {
	UserID         int    `json:"user_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	MiddleName     string `json:"middle_name"`
	PassportNumber string `json:"passport_number"`
}

type UserProfileResponse struct {
	ID         int    `json:"id"`
	Phone      string `json:"phone"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	IsVerified bool   `json:"is_verified"`
}

type UserBalanceResponse struct {
	Balance      float64 `json:"balance"`
	BonusBalance float64 `json:"bonus_balance"`
}
