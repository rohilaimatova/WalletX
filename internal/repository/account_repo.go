package repository

import (
	"WalletX/models"
	"WalletX/pkg/errs"
	"WalletX/pkg/logger"
	"context"
	"database/sql"
	"time"
)

type AccountRepository interface {
	CreateAccount(account models.Account) (models.Account, error)
	GetByID(ctx context.Context, id int) (*models.Account, error)
	DecreaseBalance(ctx context.Context, id int, amount float64) error
	IncreaseBalance(ctx context.Context, id int, amount float64) error
	GetByUserID(ctx context.Context, userID int) (models.Account, error)
	GetByPhone(ctx context.Context, phone string) (*models.Account, error)
	GetTransactions(ctx context.Context, accountID int, start, end time.Time) ([]models.TransactionHistory, error)
}

type accountRepo struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) GetByPhone(ctx context.Context, phone string) (*models.Account, error) {
	var acc models.Account
	query := `
		SELECT a.id, a.user_id, a.balance, a.bonus_balance, a.created_at, a.updated_at
		FROM accounts a
		JOIN users u ON a.user_id = u.id
		WHERE u.phone = $1
	`
	err := r.db.QueryRowContext(ctx, query, phone).
		Scan(&acc.ID, &acc.UserID, &acc.Balance, &acc.BonusBalance, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn.Printf("[AccountRepository] Account not found by phone=%s", phone)
			return nil, errs.ErrAccountNotFound
		}
		logger.Error.Printf("[AccountRepository] GetByPhone DB error: %v", err)
		return nil, errs.ErrInternal
	}

	logger.Info.Printf("[AccountRepository] Found account by phone %s: %+v", phone, acc)
	return &acc, nil
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
		return models.Account{}, errs.ErrInternal
	}
	logger.Info.Printf("[CreateAccount] success: userID=%d, account=%v", account.UserID, account) // Если успех
	return account, nil
}

func getTx(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(struct{}{}).(*sql.Tx)
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
			logger.Warn.Printf("[AccountRepository] Account not found: id=%d", id)
			return nil, errs.ErrAccountNotFound
		}
		logger.Error.Printf("[AccountRepository] GetByID DB error: %v", err)
		return nil, errs.ErrInternal
	}

	return &account, nil
}
func (r *accountRepo) GetByUserID(ctx context.Context, userID int) (models.Account, error) {
	var account models.Account
	err := r.db.QueryRowContext(ctx,
		"SELECT id, user_id, balance, bonus_balance, created_at, updated_at FROM accounts WHERE user_id = $1",
		userID,
	).Scan(&account.ID, &account.UserID, &account.Balance, &account.BonusBalance, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn.Printf("[AccountRepository] Account not found for userID=%d", userID)
			return models.Account{}, errs.ErrAccountNotFound
		}
		logger.Error.Printf("[AccountRepository] GetByUserID DB error: %v", err)
		return models.Account{}, errs.ErrInternal
	}

	return account, nil
}

func (r *accountRepo) DecreaseBalance(ctx context.Context, id int, amount float64) error {
	tx := getTx(ctx)

	var balance float64
	row := r.db.QueryRow("SELECT balance FROM accounts WHERE id = $1", id)
	if err := row.Scan(&balance); err != nil {
		logger.Error.Printf(
			"[AccountRepository] Failed to fetch balance for accountID=%d: %v",
			id, err,
		)
		if err == sql.ErrNoRows {
			return errs.ErrAccountNotFound
		}
		return errs.ErrInternal
	}

	if balance < amount {
		logger.Warn.Printf("[AccountRepository] Insufficient balance for account ID %d, requested: %.2f, available: %.2f", id, amount, balance)
		return errs.ErrInsufficientBalance
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
		return errs.ErrInternal
	}

	rows, _ := exec.RowsAffected()
	if rows == 0 {
		logger.Warn.Printf("[AccountRepository] DecreaseBalance failed: no rows affected for account ID %d", id)
		return errs.ErrAccountNotFound
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
		return errs.ErrInternal
	}

	rows, _ := exec.RowsAffected()
	if rows == 0 {
		return errs.ErrAccountNotFound
	}

	return nil
}
