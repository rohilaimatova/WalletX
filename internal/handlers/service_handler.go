package handlers

import (
	"WalletX/internal/service"
	"WalletX/pkg/logger"
	"WalletX/pkg/respond"
	"net/http"
)

type ServicesHandler struct {
	Service *service.ServicesService
}

func NewServicesHandler(s *service.ServicesService) *ServicesHandler {
	return &ServicesHandler{Service: s}
}

// GetAllServices godoc
// @Summary Get all services
// @Description Returns a list of all available services
// @Tags services
// @Accept json
// @Produce json
// @Success 200 {array} models.Services
// @Failure 500 {object} models.ErrorResponse "internal server error"
// @Security BearerAuth
// @Router /api/services [get]
func (h *ServicesHandler) GetAllServices(w http.ResponseWriter, r *http.Request) {
	logger.Info.Printf("Received request to fetch all services: %s %s", r.Method, r.URL.Path)
	services, err := h.Service.GetAllServices()

	if err != nil {
		logger.Error.Printf("Error occurred while fetching services: %v", err)
		respond.Error(w, http.StatusInternalServerError, "failed to load services", err)
		return
	}
	logger.Info.Printf("Successfully fetched %d services", len(services))

	respond.JSON(w, http.StatusOK, services)
}
