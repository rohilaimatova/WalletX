package service

import (
	"WalletX/internal/repository"
	"WalletX/models"
	"WalletX/pkg/errs"
	"WalletX/pkg/logger"
)

type UserProfileService struct {
	repo repository.UserProfileRepository
}

func NewUserProfileService(r repository.UserProfileRepository) *UserProfileService {
	return &UserProfileService{
		repo: r,
	}
}

func (s *UserProfileService) GetProfileByID(id int) (*models.UserProfileResponse, error) {
	logger.Info.Printf("[UserProfileService] GetProfileByID called with id=%d", id)

	user, err := s.repo.GetProfileByID(id)
	if err != nil {
		logger.Error.Printf("[UserProfileService] Failed to get user by ID=%d: %v", id, err)
		return nil, errs.ErrInternal
	}

	logger.Info.Printf("[UserProfileService] Successfully retrieved profile for userID=%d: %+v", id, user)
	return &user, nil
}

func (s *UserProfileService) GetBalanceByUserID(id int) (*models.UserBalanceResponse, error) {
	logger.Info.Printf("[UserProfileService] GetBalanceByUserID called with userID=%d", id)

	balance, err := s.repo.GetBalanceByUserID(id)
	if err != nil {
		logger.Error.Printf("[UserProfileService] Failed to get balance for userID=%d: %v", id, err)
		return nil, errs.ErrInternal
	}

	logger.Info.Printf("[UserProfileService] Successfully retrieved balance for userID=%d: %+v", id, balance)
	return &balance, nil
}
