package handlers

import (
	"WalletX/internal/service"
	"WalletX/pkg/logger"
	"WalletX/respond"
	"net/http"
)

type ServiceHandler struct {
	Service *service.ServiceService
}

func NewServiceHandler(s *service.ServiceService) *ServiceHandler {
	return &ServiceHandler{Service: s}
}

func (h *ServiceHandler) GetAllServices(w http.ResponseWriter, r *http.Request) {
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
