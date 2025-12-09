package service

import (
	"WalletX/internal/repository"
	"WalletX/models"
	"WalletX/pkg/errs"
	"WalletX/pkg/logger"
	"time"
)

type CardService struct {
	Repo repository.CardRepository
}

func NewCardService(repo repository.CardRepository) *CardService {
	return &CardService{Repo: repo}
}

func (s *CardService) CreateCard(userID, accountID int, cardType, maskedNumber, expiryDate string) (models.Card, error) {
	if maskedNumber == "" || cardType == "" || expiryDate == "" {
		logger.Warn.Printf("CreateCard: missing required fields for userID=%d", userID)
		return models.Card{}, errs.ErrRequiredFields
	}

	card := models.Card{
		UserID:       userID,
		AccountID:    accountID,
		MaskedNumber: maskedNumber, // маска вместо полного номера
		CardType:     cardType,
		ExpiryDate:   expiryDate,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	created, err := s.Repo.CreateCard(card)
	if err != nil {
		logger.Error.Printf("CreateCard: failed to create card for userID=%d, err=%v", userID, err)
		return models.Card{}, err
	}

	logger.Info.Printf("CreateCard: success, cardID=%d for userID=%d", created.ID, userID)
	return created, nil
}

func (s *CardService) GetCardsByUser(userID int) ([]models.Card, error) {
	cards, err := s.Repo.GetByUserID(userID)
	if err != nil {
		logger.Warn.Printf("GetCardsByUser: failed for userID=%d, err=%v", userID, err)
		return nil, err
	}
	return cards, nil
}

func (s *CardService) DeactivateCard(cardID int) error {
	err := s.Repo.DeactivateCard(cardID)
	if err != nil {
		logger.Warn.Printf("DeactivateCard: failed for cardID=%d, err=%v", cardID, err)
		return err
	}
	logger.Info.Printf("DeactivateCard: success for cardID=%d", cardID)
	return nil
}
