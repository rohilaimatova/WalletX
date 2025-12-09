package repository

import (
	"WalletX/models"
	"WalletX/pkg/logger"
	"database/sql"
)

type ServiceRepository interface {
	GetAll() ([]models.Service, error)
}

type serviceRepo struct {
	db *sql.DB
}

func NewServiceRepo(db *sql.DB) ServiceRepository {
	return &serviceRepo{db: db}
}

func (r *serviceRepo) GetAll() ([]models.Service, error) {
	logger.Info.Println("Fetching all services from the database")

	rows, err := r.db.Query(`SELECT id, name, description, created_at, updated_at FROM services`)
	if err != nil {
		logger.Error.Printf("Error occurred while fetching services: %v", err)
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			logger.Error.Printf("Error occurred while scanning row: %v", err)
		}

		services = append(services, s)

	}
	logger.Info.Printf("Successfully fetched %d services", len(services))
	return services, nil
}
