package repository

import (
	"WalletX/models"
	"database/sql"
)

type UserRepository interface {
	CreateUser(user models.User) (models.User, error)
	GetByPhone(phone string) (models.User, error)
	UpdatePassword(userID int, hashedPassword string) error
	UpdateVerification(userID int, firstName, lastName, middleName, passport string) error
	GetByID(userID int) (models.User, error)
	IncrementPasswordAttempts(userID int) error
	ResetPasswordAttempts(userID int) error
	BlockUser(userID int) error
}

type PostgresUserRepo struct {
	DB *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{DB: db}
}

func (r *PostgresUserRepo) CreateUser(user models.User) (models.User, error) {
	err := r.DB.QueryRow(
		`INSERT INTO users (phone, is_verified) VALUES ($1, $2) RETURNING id, phone, is_verified`,
		user.Phone, user.IsVerified,
	).Scan(&user.ID, &user.Phone, &user.IsVerified)
	if err != nil {
		return models.User{}, translateError(err)
	}
	return user, nil
}

func (r *PostgresUserRepo) GetByPhone(phone string) (models.User, error) {
	var user models.User
	err := r.DB.QueryRow(
		`SELECT id, phone, password, device_id, password_attempts, is_verified, first_name, last_name, middle_name, passport_number 
    FROM users WHERE phone=$1`,
		phone,
	).Scan(&user.ID, &user.Phone, &user.Password, &user.DeviceID, &user.PasswordAttempts,
		&user.IsVerified, &user.FirstName, &user.LastName, &user.MiddleName, &user.PassportNumber)
	return user, translateError(err)
}
func (r *PostgresUserRepo) UpdatePassword(userID int, hashedPassword string) error {
	_, err := r.DB.Exec(
		`UPDATE users SET password=$1 WHERE id=$2`,
		hashedPassword, userID,
	)
	return translateError(err)
}
func (r *PostgresUserRepo) IncrementPasswordAttempts(userID int) error {
	_, err := r.DB.Exec(
		`UPDATE users SET password_attempts = password_attempts + 1 WHERE id=$1`,
		userID,
	)
	return translateError(err)
}

func (r *PostgresUserRepo) ResetPasswordAttempts(userID int) error {
	_, err := r.DB.Exec(
		`UPDATE users SET password_attempts = 0 WHERE id=$1`,
		userID,
	)
	return translateError(err)
}

func (r *PostgresUserRepo) BlockUser(userID int) error {
	_, err := r.DB.Exec(
		`UPDATE users SET device_id = true WHERE id=$1`,
		userID,
	)
	return translateError(err)
}

func (r *PostgresUserRepo) UpdateVerification(userID int, firstName, lastName, middleName, passport_number string) error {
	_, err := r.DB.Exec(
		`UPDATE users 
		 SET first_name=$1, last_name=$2, middle_name=$3, passport_number=$4, is_verified=true 
		 WHERE id=$5`,
		firstName, lastName, middleName, passport_number, userID,
	)
	return translateError(err)
}

func (r *PostgresUserRepo) GetByID(userID int) (models.User, error) {
	var user models.User
	err := r.DB.QueryRow(
		`SELECT id, phone, password, is_verified, first_name, last_name, middle_name, passport_number
		FROM users WHERE id=$1`,
		userID,
	).Scan(&user.ID, &user.Phone, &user.Password, &user.IsVerified,
		&user.FirstName, &user.LastName, &user.MiddleName, &user.PassportNumber)
	return user, translateError(err)
}
