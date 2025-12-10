package repository

import (
	"WalletX/models"
	"WalletX/pkg/logger"
	"context"
	"database/sql"
	"errors"
)

type AccountRepository interface {
	CreateAccount(account models.Account) (models.Account, error)
	GetByID(ctx context.Context, id int) (*models.Account, error)
	DecreaseBalance(ctx context.Context, id int, amount float64) error
	IncreaseBalance(ctx context.Context, id int, amount float64) error
}

type accountRepo struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) CreateAccount(account models.Account) (models.Account, error) {
	query := `
		INSERT INTO accounts (user_id, balance, bonus_balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	row := r.db.QueryRow(query, account.UserID, account.Balance, account.BonusBalance, account.CreatedAt, account.UpdatedAt)
	err := row.Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		logger.Warn.Printf("[CreateAccount] failed: userID=%d, err=%v", account.UserID, err) // Если ошибка
		return models.Account{}, err
	}
	logger.Info.Printf("[CreateAccount] success: userID=%d, account=%v", account.UserID, account) // Если успех
	return account, nil
}

func getTx(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(struct{}{}).(*sql.Tx) // если используешь другой ключ, меняем
	return tx
}

func (r *accountRepo) GetByID(ctx context.Context, id int) (*models.Account, error) {
	tx := getTx(ctx)

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow("SELECT id, user_id, balance, bonus_balance, created_at, updated_at FROM accounts WHERE id = $1", id)
	} else {
		row = r.db.QueryRow("SELECT id, user_id, balance, bonus_balance, created_at, updated_at FROM accounts WHERE id = $1", id)
	}

	var account models.Account
	err := row.Scan(&account.ID, &account.UserID, &account.Balance, &account.BonusBalance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если аккаунт не найден
			logger.Warn.Printf("[AccountRepository] GetByID: no account found for ID=%d", id)
			return nil, err
		}
		logger.Error.Printf("[AccountRepository] GetByID error: %v", err)
		return nil, err
	}

	return &account, nil
}

func (r *accountRepo) DecreaseBalance(ctx context.Context, id int, amount float64) error {
	tx := getTx(ctx)

	var balance float64
	row := r.db.QueryRow("SELECT balance FROM accounts WHERE id = $1", id)
	if err := row.Scan(&balance); err != nil {
		logger.Error.Printf("[AccountRepository] Error fetching balance for account ID %d: %v", id, err)
		return err
	}
	logger.Info.Printf("[AccountRepository] Account ID %d current balance: %.2f", id, balance)

	if balance < amount {
		logger.Warn.Printf("[AccountRepository] Insufficient balance for account ID %d, requested: %.2f, available: %.2f", id, amount, balance)
		return errors.New("insufficient balance")
	}

	var exec sql.Result
	var err error
	if tx != nil {
		exec, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1", amount, id)
	} else {
		exec, err = r.db.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1", amount, id)
	}

	if err != nil {
		logger.Error.Printf("[AccountRepository] DecreaseBalance error: %v", err)
		return err
	}

	rows, _ := exec.RowsAffected()
	if rows == 0 {
		// Если строка не была обновлена, значит либо аккаунт не найден, либо недостаточно средств
		logger.Warn.Printf("[AccountRepository] DecreaseBalance failed: no rows affected for account ID %d", id)
		return errors.New("insufficient balance or account not found")
	}

	logger.Info.Printf("[AccountRepository] Successfully decreased balance for account ID %d, amount: %.2f", id, amount)

	return nil
}

func (r *accountRepo) IncreaseBalance(ctx context.Context, id int, amount float64) error {
	tx := getTx(ctx)

	var exec sql.Result
	var err error
	if tx != nil {
		exec, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, id)
	} else {
		exec, err = r.db.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, id)
	}

	if err != nil {
		logger.Error.Printf("[AccountRepository] IncreaseBalance error: %v", err)
		return err
	}

	rows, _ := exec.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
