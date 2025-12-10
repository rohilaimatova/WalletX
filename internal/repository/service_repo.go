package repository

import (
	"WalletX/models"
	"WalletX/pkg/logger"
	"database/sql"
	"fmt"
)

type ServiceRepository interface {
	GetAll() ([]models.Service, error)
	GetServiceIDByType(serviceType string) (int, error)
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
func (r *serviceRepo) GetByID(id int) (*models.Service, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM services WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var s models.Service
	err := row.Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service with id %d not found", id)
		}
		return nil, err
	}

	return &s, nil
}
func (r *serviceRepo) GetServiceIDByType(serviceType string) (int, error) {
	var id int
	query := `SELECT id FROM services WHERE name = $1`
	err := r.db.QueryRow(query, serviceType).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn.Printf("[ServiceRepository] Service type not found: %s", serviceType)
			return 0, fmt.Errorf("service type %s not found", serviceType)
		}
		logger.Error.Printf("[ServiceRepository] GetServiceIDByType error: %v", err)
		return 0, err
	}
	return id, nil
}
