package service

import (
	"WalletX/internal/handlers/transaction"
	"WalletX/internal/repository"
	"WalletX/models"
	"WalletX/pkg/logger"
	"context"
	"errors"
	"time"
)

type AccountService struct {
	Repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{Repo: repo}
}

func (s *AccountService) CreateAccountForUser(userID int) (models.Account, error) {
	account := models.Account{
		UserID:       userID,
		Balance:      0,
		BonusBalance: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	created, err := s.Repo.CreateAccount(account)
	if err != nil {
		logger.Error.Printf("CreateAccountForUser: failed for userID=%d, err=%v", userID, err)
		return models.Account{}, err
	}

	logger.Info.Printf("CreateAccountForUser: success, accountID=%d for userID=%d", created.ID, userID)
	return created, nil
}

type PaymentService struct {
	AccountRepo     repository.AccountRepository
	TransactionRepo repository.TransactionRepository
	ServiceRepo     repository.ServiceRepository
	TM              transaction.TransactionManager
}

func NewPaymentService(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository, serviceRepo repository.ServiceRepository, tm transaction.TransactionManager) *PaymentService {
	return &PaymentService{
		AccountRepo:     accountRepo,
		TransactionRepo: transactionRepo,
		ServiceRepo:     serviceRepo,
		TM:              tm,
	}
}

func (s *PaymentService) Pay(ctx context.Context, fromID, toID int, amount float64, transactionType string) error {
	return s.TM.WithinTransaction(ctx, func(txCtx context.Context) error {

		from, err := s.AccountRepo.GetByID(txCtx, fromID)
		if err != nil {
			logger.Warn.Printf("[PaymentService] account_from not found: %d", fromID)
			return errors.New("account_from not found")
		}

		to, err := s.AccountRepo.GetByID(txCtx, toID)
		if err != nil {
			logger.Warn.Printf("[PaymentService] account_to not found: %d", toID)
			return errors.New("account_to not found")
		}
		logger.Info.Printf("[PaymentService] paying from %d to %d with amount %.2f", from.ID, to.ID, amount)

		if from.Balance < amount {
			logger.Warn.Printf("[PaymentService] insufficient balance: have=%.2f need=%.2f", from.Balance, amount)
			return errors.New("insufficient balance")
		}

		err = s.AccountRepo.DecreaseBalance(txCtx, fromID, amount)
		if err != nil {
			logger.Error.Printf("[PaymentService] DecreaseBalance error: %v", err)
			return err
		}

		err = s.AccountRepo.IncreaseBalance(txCtx, toID, amount)
		if err != nil {
			logger.Error.Printf("[PaymentService] IncreaseBalance error: %v", err)
			return err
		}

		transaction := models.Transaction{
			AccountFrom: fromID,
			AccountTo:   toID,
			Amount:      amount,
			Type:        transactionType,
			CreatedAt:   time.Now(),
		}

		_, err = s.TransactionRepo.CreateTransaction(transaction)
		if err != nil {
			logger.Error.Printf("[PaymentService] Failed to create transaction: %v", err)
			return err
		}

		logger.Info.Printf("[PaymentService] SUCCESS transfer from=%d to=%d amount=%.2f type=%s", fromID, toID, amount, transactionType)
		return nil
	})
}
