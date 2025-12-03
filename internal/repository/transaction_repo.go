package repository

import (
	"WalletX/models"
	"database/sql"
)

type TransactionRepository interface {
	SaveTransaction(tx *models.Transaction) error
	GetTransactionsByAccountID(accountID int) ([]models.Transaction, error)
}

type PostgresTransactionRepo struct {
	DB *sql.DB
}

func NewPostgresTransactionRepo(db *sql.DB) *PostgresTransactionRepo {
	return &PostgresTransactionRepo{DB: db}
}

// Сохраняем транзакцию
func (r *PostgresTransactionRepo) SaveTransaction(tx *models.Transaction) error {
	return r.DB.QueryRow(
		"INSERT INTO transactions (account_from, account_to, type, amount) VALUES ($1, $2, $3, $4) RETURNING id",
		tx.AccountFrom, tx.AccountTo, tx.Type, tx.Amount,
	).Scan(&tx.ID)
}

// Получаем все транзакции, где аккаунт был либо отправителем, либо получателем
func (r *PostgresTransactionRepo) GetTransactionsByAccountID(accountID int) ([]models.Transaction, error) {
	rows, err := r.DB.Query(
		`SELECT id, account_from, account_to, type, amount, created_at
		 FROM transactions 
		 WHERE account_from=$1 OR account_to=$1
		 ORDER BY created_at DESC`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.AccountFrom, &tx.AccountTo, &tx.Type, &tx.Amount, &tx.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, tx)
	}
	return list, nil
}
