package repository

import (
	"WalletX/models"
	"database/sql"
)

type AccountRepository interface {
	CreateAccount(userID int) (models.Account, error)
	GetAccountByID(id int) (models.Account, error)
	UpdateBalance(id int, balance float64) error
}

type PostgresAccountRepo struct {
	DB *sql.DB
}

func NewPostgresAccountRepo(db *sql.DB) *PostgresAccountRepo {
	return &PostgresAccountRepo{DB: db}
}

// Создаём аккаунт для конкретного пользователя
func (r *PostgresAccountRepo) CreateAccount(userID int) (models.Account, error) {
	var acc models.Account

	err := r.DB.QueryRow(
		"INSERT INTO accounts (user_id, balance) VALUES ($1, 0) RETURNING id, user_id, balance",
		userID,
	).Scan(&acc.ID, &acc.UserID, &acc.Balance)

	return acc, err
}

func (r *PostgresAccountRepo) GetAccountByID(id int) (models.Account, error) {
	var acc models.Account

	err := r.DB.QueryRow(
		"SELECT id, user_id, balance FROM accounts WHERE id=$1",
		id,
	).Scan(&acc.ID, &acc.UserID, &acc.Balance)

	return acc, err
}

func (r *PostgresAccountRepo) UpdateBalance(id int, balance float64) error {
	_, err := r.DB.Exec(
		"UPDATE accounts SET balance=$1 WHERE id=$2",
		balance, id,
	)
	return err
}
