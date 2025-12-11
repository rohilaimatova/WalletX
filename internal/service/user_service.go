package service

import (
	"WalletX/internal/repository"
	"WalletX/models"
	"fmt"
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
	user, err := s.repo.GetProfileByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (s *UserProfileService) GetBalanceByUserID(id int) (*models.UserBalanceResponse, error) {
	balance, err := s.repo.GetBalanceByUserID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return &balance, nil
}
