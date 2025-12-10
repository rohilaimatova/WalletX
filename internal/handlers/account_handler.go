// internal/handlers/account_handler.go

package handlers

import (
	"WalletX/internal/service"
	"WalletX/models"
	"WalletX/pkg/logger"
	"WalletX/respond"
	"encoding/json"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	Payment *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *AccountHandler {
	return &AccountHandler{Payment: paymentService}
}

func (h *AccountHandler) PayForService(w http.ResponseWriter, r *http.Request) {
	logger.Info.Printf("[PayForService] Incoming request: %s %s", r.Method, r.URL.Path)

	var req models.PayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error.Printf("[PayForService] Failed to decode request: %v", err)
		respond.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	logger.Info.Printf("[PayForService] Request payload: %+v", req)

	fromID, err := strconv.Atoi(req.Account)
	if err != nil {
		logger.Warn.Printf("[PayForService] Invalid account id: %s", req.Account)
		respond.Error(w, http.StatusBadRequest, "invalid account id", err)
		return
	}
	logger.Info.Printf("[PayForService] fromID parsed: %d", fromID)

	toID, err := h.Payment.ServiceRepo.GetServiceIDByType(req.ServiceType)
	if err != nil {
		logger.Warn.Printf("[PayForService] Invalid service type: %s, error: %v", req.ServiceType, err)
		respond.Error(w, http.StatusBadRequest, "invalid service type", err)
		return
	}
	logger.Info.Printf("[PayForService] toID resolved: %d", toID)

	transactionType := req.ServiceType

	err = h.Payment.Pay(r.Context(), fromID, toID, req.Amount, transactionType) // Передаем тип транзакции
	if err != nil {
		logger.Error.Printf("[PayForService] Payment failed from=%d to=%d amount=%.2f: %v", fromID, toID, req.Amount, err)
		respond.Error(w, http.StatusBadRequest, "payment failed", err)
		return
	}

	logger.Info.Printf("[PayForService] Payment success from=%d to=%d amount=%.2f type=%s", fromID, toID, req.Amount, transactionType)
	respond.JSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "payment completed",
	})
}
