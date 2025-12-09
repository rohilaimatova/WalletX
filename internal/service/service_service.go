package service

import (
	"WalletX/internal/repository"
	"WalletX/models"
	"WalletX/pkg/logger" // Импорт кастомного логгера
)

type ServiceService struct {
	Repo repository.ServiceRepository
}

func NewServiceService(repo repository.ServiceRepository) *ServiceService {
	return &ServiceService{Repo: repo}
}

func (s *ServiceService) GetAllServices() ([]models.Service, error) {
	logger.Info.Println("Fetching all services")

	services, err := s.Repo.GetAll()
	if err != nil {
		logger.Error.Printf("Error occurred while fetching services from repository: %v", err)
		return nil, err
	}

	logger.Info.Printf("Successfully fetched %d services", len(services))
	return services, nil
}
