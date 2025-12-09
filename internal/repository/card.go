package repository

import (
	"WalletX/models"
	"database/sql"
)

type CardRepository interface {
	CreateCard(card models.Card) (models.Card, error)
	GetByUserID(userID int) ([]models.Card, error)
	DeactivateCard(cardID int) error
}

type PostgresCardRepo struct {
	DB *sql.DB
}

func NewPostgresCardRepo(db *sql.DB) *PostgresCardRepo {
	return &PostgresCardRepo{DB: db}
}

func (r *PostgresCardRepo) CreateCard(card models.Card) (models.Card, error) {
	err := r.DB.QueryRow(
		`INSERT INTO cards (user_id, account_id, masked_number, card_type, expiry_date, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		card.UserID, card.AccountID, card.MaskedNumber, card.CardType, card.ExpiryDate,
		card.IsActive, card.CreatedAt, card.UpdatedAt,
	).Scan(&card.ID)
	if err != nil {
		return models.Card{}, err
	}
	return card, nil
}

func (r *PostgresCardRepo) GetByUserID(userID int) ([]models.Card, error) {
	rows, err := r.DB.Query(
		`SELECT id, user_id, account_id, masked_number, card_type, expiry_date, is_active, created_at, updated_at 
		 FROM cards WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var c models.Card
		if err := rows.Scan(
			&c.ID, &c.UserID, &c.AccountID, &c.MaskedNumber, &c.CardType, &c.ExpiryDate,
			&c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func (r *PostgresCardRepo) DeactivateCard(cardID int) error {
	res, err := r.DB.Exec(`UPDATE cards SET is_active=false, updated_at=NOW() WHERE id=$1`, cardID)
	if err != nil {
		return err
	}

	// Проверяем, что реально обновили запись
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
