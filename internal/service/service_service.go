package service

import (
	"WalletX/internal/repository"
	"WalletX/models"
	"WalletX/pkg/logger"
)

type ServicesService struct {
	Repo repository.ServicesRepository
}

func NewServicesService(repo repository.ServicesRepository) *ServicesService {
	return &ServicesService{
		Repo: repo,
	}
}

func (s *ServicesService) GetAllServices() ([]models.Services, error) {
	logger.Info.Println("Fetching all services")
	services, err := s.Repo.GetAll()
	if err != nil {
		logger.Error.Printf("Error occurred while fetching services from repository: %v", err)
		return nil, err
	}
	logger.Info.Printf("Successfully fetched %d services", len(services))
	return services, nil
}
