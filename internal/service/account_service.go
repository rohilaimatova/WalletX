package service

import (
	"WalletX/internal/repository"
	"WalletX/models"
	"WalletX/pkg/logger"
	"errors"
)

type AccountService struct {
	Repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{Repo: repo}
}

func (s *AccountService) CreateAccount(req models.CreateAccountRequest) (models.Account, error) {
	if req.UserID == 0 {
		logger.Warn.Printf("CreateAccount: invalid input, UserID not provided")
		return models.Account{}, errors.New("user ID must be provided")
	}

	acc, err := s.Repo.CreateAccount(req.UserID)
	if err != nil {
		logger.Error.Printf("CreateAccount: failed to create account for UserID %d: %v", req.UserID, err)
		return models.Account{}, err
	}

	logger.Info.Printf("CreateAccount: successfully created account ID %d for UserID %d", acc.ID, req.UserID)
	return acc, nil
}

func (s *AccountService) GetBalance(accountID int) (float64, error) {
	acc, err := s.Repo.GetAccountByID(accountID)
	if err != nil {
		logger.Error.Printf("GetBalance: failed to get account ID %d: %v", accountID, err)
		return 0, err
	}

	logger.Info.Printf("GetBalance: account ID %d has balance %.2f", accountID, acc.Balance)
	return acc.Balance, nil
}
