package repository

import (
	"WalletX/models"
	"database/sql"
	"fmt"
)

type UserProfileRepository interface {
	GetProfileByID(id int) (models.UserProfileResponse, error)
	GetBalanceByUserID(userID int) (models.UserBalanceResponse, error)
}

type userProfileRepo struct {
	db *sql.DB
}

func NewUserProfileRepository(db *sql.DB) UserProfileRepository {
	return &userProfileRepo{
		db: db,
	}
}

func (r *userProfileRepo) GetProfileByID(id int) (models.UserProfileResponse, error) {
	var user models.UserProfileResponse

	query := `
        SELECT id, phone, first_name, last_name, middle_name, is_verified
        FROM users
        WHERE id = $1
        LIMIT 1
    `
	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&user.ID,
		&user.Phone,
		&user.FirstName,
		&user.LastName,
		&user.MiddleName,
		&user.IsVerified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return user, fmt.Errorf("failed to scan user: %w", err)
	}

	return user, nil
}

func (r *userProfileRepo) GetBalanceByUserID(userID int) (models.UserBalanceResponse, error) {
	var balance models.UserBalanceResponse

	query := `
        SELECT balance, bonus_balance
        FROM accounts
        WHERE user_id = $1
        LIMIT 1
    `
	row := r.db.QueryRow(query, userID)
	err := row.Scan(&balance.Balance, &balance.BonusBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return balance, nil
		}
		return balance, fmt.Errorf("failed to scan balance: %w", err)
	}

	return balance, nil
}
