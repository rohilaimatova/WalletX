package repository

import (
	"WalletX/models"
	"database/sql"
	"time"
)

type AccountRepository interface {
	CreateAccount(account models.Account) (models.Account, error)
	GetByUserID(userID int) (models.Account, error)
	UpdateBalance(userID int, balance, bonus float64) error
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
		return models.Account{}, err
	}
	return account, nil
}

func (r *accountRepo) GetByUserID(userID int) (models.Account, error) {
	var account models.Account
	query := `SELECT id, user_id, balance, bonus_balance, created_at, updated_at FROM accounts WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.BonusBalance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return models.Account{}, err
	}
	return account, nil
}

func (r *accountRepo) UpdateBalance(userID int, balance, bonus float64) error {
	query := `UPDATE accounts SET balance = $1, bonus_balance = $2, updated_at = $3 WHERE user_id = $4`
	_, err := r.db.Exec(query, balance, bonus, time.Now(), userID)
	return err
}

/*package repository

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
*/
