package models

type User struct {
	ID               int     `json:"id"`
	Phone            string  `json:"phone"`
	Password         string  `json:"-"`
	FirstName        *string `json:"first_name,omitempty"`
	LastName         *string `json:"last_name,omitempty"`
	MiddleName       *string `json:"middle_name,omitempty"`
	PassportNumber   *string `json:"passport_number,omitempty"`
	IsVerified       bool    `json:"is_verified"`
	DeviceID         bool    `json:"device_id,omitempty"`
	PasswordAttempts int     `json:"password_attempts" db:"password_attempts"`
}
type SetPasswordRequest struct {
	UserID   int    `json:"user_id"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Phone    string `json:"phone" example:"+992931753756"`
	Password string `json:"password" example:"12345678"`
}
type VerifyIdentityRequest struct {
	UserID         int    `json:"user_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	MiddleName     string `json:"middle_name"`
	PassportNumber string `json:"passport_number"`
}

type UserProfileResponse struct {
	ID         int    `json:"id" example:"1"`
	Phone      string `json:"phone" example:"+992931753756"`
	FirstName  string `json:"first_name" example: "Ali"`
	LastName   string `json:"last_name" example: "Bob"`
	MiddleName string `json:"middle_name" example: "Bob"`
	IsVerified bool   `json:"is_verified" example:"true"`
}

type UserBalanceResponse struct {
	Balance      float64 `json:"balance" example:"100"`
	BonusBalance float64 `json:"bonus_balance" example:"100"`
}
