package errs

import "errors"

var (
	ErrInvalidPhone        = errors.New("invalid phone number format")
	ErrUserExists          = errors.New("user already exists")
	ErrWeakPassword        = errors.New("password must be exactly 8 characters long")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrUserNotFound        = errors.New("user not found")
	ErrRequiredFields      = errors.New("all required fields must be filled")
	ErrInternal            = errors.New("internal error")
	ErrValidationFailed    = errors.New("validation failed")
	ErrUserBlocked         = errors.New("user is blocked")
	ErrWrongPassword       = errors.New("wrong password")
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrSelfTransfer        = errors.New("cannot transfer to yourself")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInsufficientBalance = errors.New("insufficient balance or account not found")
	ErrAccountNotFound     = errors.New("account not found")
)
