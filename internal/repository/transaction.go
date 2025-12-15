package repository

import (
	"WalletX/models"
	"WalletX/pkg/errs"
	"WalletX/pkg/logger"
	"context"
	"database/sql"
	"time"
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

func (r *accountRepo) GetTransactions(ctx context.Context, accountID int, start, end time.Time) ([]models.TransactionHistory, error) {
	logger.Info.Printf(
		"[AccountRepository] Fetching transactions for accountID=%d, period=%s - %s",
		accountID, start, end,
	)

	query := `
		SELECT
			t.account_to,
			t.amount,
			t.type,
			t.created_at,
			u.phone
		FROM transactions t
		LEFT JOIN accounts a ON a.id = t.account_to
		LEFT JOIN users u ON u.id = a.user_id
		WHERE t.account_from = $1
		  AND t.created_at BETWEEN $2 AND $3
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, start, end)
	if err != nil {
		logger.Error.Printf("[AccountRepository] Failed to fetch transactions: %v", err)
		return nil, errs.ErrInternal
	}
	defer rows.Close()

	txs := make([]models.TransactionHistory, 0)

	for rows.Next() {
		var t models.TransactionHistory
		var phone sql.NullString

		if err := rows.Scan(
			&t.AccountTo,
			&t.Amount,
			&t.Type,
			&t.CreatedAt,
			&phone,
		); err != nil {
			logger.Error.Printf("[AccountRepository] Scan error: %v", err)
			continue
		}

		if t.Type == "transfer" && phone.Valid {
			t.ToPhone = &phone.String
		} else {
			t.ToPhone = nil
		}

		txs = append(txs, t)
	}

	logger.Info.Printf(
		"[AccountRepository] Fetched %d transactions for accountID=%d",
		len(txs), accountID,
	)

	return txs, nil
}
