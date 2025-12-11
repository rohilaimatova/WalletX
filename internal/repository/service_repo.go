package repository

import (
	"WalletX/models"
	"WalletX/pkg/logger"
	"database/sql"
	"fmt"
)

type ServicesRepository interface {
	GetAll() ([]models.Services, error)
	GetServiceIDByType(serviceType string) (int, error)
	GetByID(id int) (*models.Services, error)
}

type servicesRepo struct {
	db *sql.DB
}

func NewServicesRepo(db *sql.DB) ServicesRepository {
	return &servicesRepo{db: db}
}

func (r *servicesRepo) GetAll() ([]models.Services, error) {
	logger.Info.Println("Fetching all services from the database")

	rows, err := r.db.Query(`SELECT id, name, description, created_at, updated_at FROM services`)
	if err != nil {
		logger.Error.Printf("Error occurred while fetching services: %v", err)
		return nil, err
	}
	defer rows.Close()

	var services []models.Services
	for rows.Next() {
		var s models.Services
		err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			logger.Error.Printf("Error while scanning row: %v", err)
			continue
		}
		services = append(services, s)
	}
	logger.Info.Printf("Successfully fetched %d services", len(services))
	return services, nil
}

func (r *servicesRepo) GetByID(id int) (*models.Services, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM services WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var s models.Services
	err := row.Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service with id %d not found", id)
		}
		return nil, err
	}

	return &s, nil
}

func (r *servicesRepo) GetServiceIDByType(serviceType string) (int, error) {
	var id int
	query := `SELECT id FROM services WHERE name = $1`

	err := r.db.QueryRow(query, serviceType).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn.Printf("[ServicesRepository] Service type not found: %s", serviceType)
			return 0, fmt.Errorf("service type %s not found", serviceType)
		}
		logger.Error.Printf("[ServicesRepository] GetServiceIDByType error: %v", err)
		return 0, err
	}

	return id, nil
}
