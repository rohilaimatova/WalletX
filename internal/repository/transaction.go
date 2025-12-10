package repository

import (
	"WalletX/models"
	"WalletX/pkg/logger"
	"database/sql"
)

type TransactionRepository interface {
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
}

type transactionRepo struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	query := `
        INSERT INTO transactions (account_from, account_to, amount, type, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `
	row := r.db.QueryRow(query, transaction.AccountFrom, transaction.AccountTo, transaction.Amount, transaction.Type, transaction.CreatedAt)
	err := row.Scan(&transaction.ID, &transaction.CreatedAt)
	if err != nil {
		logger.Warn.Printf("[CreateTransaction] failed: from=%d to=%d, err=%v", transaction.AccountFrom, transaction.AccountTo, err)
		return models.Transaction{}, err
	}
	logger.Info.Printf("[CreateTransaction] success: transaction=%v", transaction)
	return transaction, nil
}
