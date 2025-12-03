package service

import (
	"WalletX/internal/repository"
	"WalletX/models"
	"WalletX/pkg/logger"
	"errors"
)

type TransactionService struct {
	AccRepo repository.AccountRepository
	TxRepo  repository.TransactionRepository
}

func NewTransactionService(a repository.AccountRepository, t repository.TransactionRepository) *TransactionService {
	return &TransactionService{AccRepo: a, TxRepo: t}
}

func (s *TransactionService) Deposit(req models.DepositRequest) (*models.Transaction, error) {
	acc, err := s.AccRepo.GetAccountByID(req.AccountID)
	if err != nil {
		logger.Error.Printf("Deposit: failed to get account ID %d: %v", req.AccountID, err)
		return nil, err
	}

	acc.Balance += req.Amount
	if err := s.AccRepo.UpdateBalance(acc.ID, acc.Balance); err != nil {
		logger.Error.Printf("Deposit: failed to update balance for account ID %d: %v", acc.ID, err)
		return nil, err
	}

	tx := &models.Transaction{
		AccountFrom: nil,
		AccountTo:   &acc.ID,
		Type:        "deposit",
		Amount:      req.Amount,
	}

	if err := s.TxRepo.SaveTransaction(tx); err != nil {
		logger.Error.Printf("Deposit: failed to save transaction for account ID %d: %v", acc.ID, err)
		return nil, err
	}

	logger.Info.Printf("Deposit: successful deposit of %.2f to account ID %d", req.Amount, acc.ID)
	return tx, nil
}

func (s *TransactionService) Withdraw(req models.WithdrawRequest) (*models.Transaction, error) {
	acc, err := s.AccRepo.GetAccountByID(req.AccountID)
	if err != nil {
		logger.Error.Printf("Withdraw: failed to get account ID %d: %v", req.AccountID, err)
		return nil, err
	}

	if acc.Balance < req.Amount {
		logger.Warn.Printf("Withdraw: insufficient balance for account ID %d", acc.ID)
		return nil, errors.New("insufficient balance")
	}

	acc.Balance -= req.Amount
	if err := s.AccRepo.UpdateBalance(acc.ID, acc.Balance); err != nil {
		logger.Error.Printf("Withdraw: failed to update balance for account ID %d: %v", acc.ID, err)
		return nil, err
	}

	tx := &models.Transaction{
		AccountFrom: &acc.ID,
		AccountTo:   nil,
		Type:        "withdraw",
		Amount:      req.Amount,
	}

	if err := s.TxRepo.SaveTransaction(tx); err != nil {
		logger.Error.Printf("Withdraw: failed to save transaction for account ID %d: %v", acc.ID, err)
		return nil, err
	}

	logger.Info.Printf("Withdraw: successful withdrawal of %.2f from account ID %d", req.Amount, acc.ID)
	return tx, nil
}

func (s *TransactionService) Transfer(req models.TransferRequest) (*models.Transaction, error) {
	from, err := s.AccRepo.GetAccountByID(req.FromID)
	if err != nil {
		logger.Error.Printf("Transfer: failed to get from account ID %d: %v", req.FromID, err)
		return nil, err
	}

	to, err := s.AccRepo.GetAccountByID(req.ToID)
	if err != nil {
		logger.Error.Printf("Transfer: failed to get to account ID %d: %v", req.ToID, err)
		return nil, err
	}

	if from.Balance < req.Amount {
		logger.Warn.Printf("Transfer: insufficient balance for account ID %d", from.ID)
		return nil, errors.New("insufficient balance")
	}

	from.Balance -= req.Amount
	to.Balance += req.Amount

	if err := s.AccRepo.UpdateBalance(from.ID, from.Balance); err != nil {
		logger.Error.Printf("Transfer: failed to update balance for from account ID %d: %v", from.ID, err)
		return nil, err
	}
	if err := s.AccRepo.UpdateBalance(to.ID, to.Balance); err != nil {
		logger.Error.Printf("Transfer: failed to update balance for to account ID %d: %v", to.ID, err)
		return nil, err
	}

	tx := &models.Transaction{
		AccountFrom: &from.ID,
		AccountTo:   &to.ID,
		Type:        "transfer",
		Amount:      req.Amount,
	}

	if err := s.TxRepo.SaveTransaction(tx); err != nil {
		logger.Error.Printf("Transfer: failed to save transaction from %d to %d: %v", from.ID, to.ID, err)
		return nil, err
	}

	logger.Info.Printf("Transfer: successful transfer of %.2f from account ID %d to account ID %d", req.Amount, from.ID, to.ID)
	return tx, nil
}

func (s *TransactionService) GetHistory(accountID int) ([]models.Transaction, error) {
	history, err := s.TxRepo.GetTransactionsByAccountID(accountID)
	if err != nil {
		logger.Error.Printf("GetHistory: failed to get history for account ID %d: %v", accountID, err)
		return nil, err
	}

	logger.Info.Printf("GetHistory: fetched %d transactions for account ID %d", len(history), accountID)
	return history, nil
}
