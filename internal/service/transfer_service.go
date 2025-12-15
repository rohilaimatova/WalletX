package service

import (
	"WalletX/internal/handlers/transaction"
	"WalletX/internal/repository"
	"WalletX/models"
	"WalletX/pkg/errs"
	"WalletX/pkg/logger"
	"context"
	"time"
)

type TransferService struct {
	AccountRepo     repository.AccountRepository
	TransactionRepo repository.TransactionRepository
	TM              transaction.TransactionManager
}

func NewTransferService(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository, tm transaction.TransactionManager) *TransferService {
	return &TransferService{
		AccountRepo:     accountRepo,
		TransactionRepo: transactionRepo,
		TM:              tm,
	}
}

func (s *TransferService) Transfer(ctx context.Context, fromAccountID, toAccountID int, amount float64) error {
	if amount <= 0 {
		logger.Warn.Printf("[TransferService] Invalid transfer amount: %.2f", amount)
		return errs.ErrInvalidAmount
	}

	return s.TM.WithinTransaction(ctx, func(txCtx context.Context) error {
		fromAcc, err := s.AccountRepo.GetByID(txCtx, fromAccountID)
		if err != nil {
			logger.Warn.Printf("[TransferService] Sender account not found: %v", err)
			return errs.ErrUserNotFound
		}

		toAcc, err := s.AccountRepo.GetByID(txCtx, toAccountID)
		if err != nil {
			logger.Warn.Printf("[TransferService] Recipient account not found: %v", err)
			return errs.ErrUserNotFound
		}

		if fromAcc.ID == toAcc.ID {
			logger.Warn.Printf("[TransferService] Attempt to transfer to self: accountID=%d", fromAcc.ID)
			return errs.ErrSelfTransfer
		}

		if fromAcc.Balance < amount {
			logger.Warn.Printf("[TransferService] Insufficient funds: fromAccountID=%d, balance=%.2f, requested=%.2f",
				fromAcc.ID, fromAcc.Balance, amount)
			return errs.ErrInsufficientFunds
		}

		if err := s.AccountRepo.DecreaseBalance(txCtx, fromAcc.ID, amount); err != nil {
			logger.Error.Printf("[TransferService] Failed to decrease balance: %v", err)
			return errs.ErrInternal
		}

		if err := s.AccountRepo.IncreaseBalance(txCtx, toAcc.ID, amount); err != nil {
			logger.Error.Printf("[TransferService] Failed to increase balance: %v", err)
			return errs.ErrInternal
		}

		// Сохраняем транзакцию
		tx := models.Transaction{
			AccountFrom: fromAcc.ID,
			AccountTo:   toAcc.ID,
			Amount:      amount,
			Type:        "transfer",
			CreatedAt:   time.Now(),
		}
		if _, err := s.TransactionRepo.CreateTransaction(tx); err != nil {
			logger.Error.Printf("[TransferService] Failed to create transaction: %v", err)
			return errs.ErrInternal
		}

		logger.Info.Printf("[TransferService] Transfer success: fromID=%d, toID=%d, amount=%.2f",
			fromAcc.ID, toAcc.ID, amount)
		return nil
	})
}
