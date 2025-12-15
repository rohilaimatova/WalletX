package handlers

import (
	"WalletX/internal/handlers/middleware"
	"WalletX/internal/service"
	"WalletX/models"
	"WalletX/pkg/logger"
	"WalletX/respond"
	"encoding/json"
	"net/http"
)

type AccountHandler struct {
	Payment *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *AccountHandler {
	return &AccountHandler{Payment: paymentService}
}

// PayForService godoc
// @Summary Pay for a service
// @Description Pay for a service like internet, mobile, etc.
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.PayRequest true "Payment request"
// @Success 200 {object} map[string]string "payment completed"
// @Failure 400 {object} models.ErrorResponse "bad request"
// @Failure 401 {object} models.ErrorResponse "unauthorized"
// @Failure 500 {object} models.ErrorResponse "internal error"
// @Router /api/pay [post]
func (h *AccountHandler) PayForService(w http.ResponseWriter, r *http.Request) {
	logger.Info.Printf("[PayForService] Incoming request: %s %s", r.Method, r.URL.Path)

	var req models.PayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error.Printf("[PayForService] Failed to decode request: %v", err)
		respond.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	userIDRaw := r.Context().Value(middleware.UserIDCtx)
	if userIDRaw == nil {
		respond.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	fromID, ok := userIDRaw.(int)
	if !ok {
		respond.Error(w, http.StatusInternalServerError, "invalid userID type", nil)
		return
	}

	toID, err := h.Payment.ServiceRepo.GetServiceIDByType(req.ServiceType)
	if err != nil {
		logger.Warn.Printf("[PayForService] Invalid service type: %s, error: %v", req.ServiceType, err)
		respond.Error(w, http.StatusBadRequest, "invalid service type", err)
		return
	}

	err = h.Payment.Pay(r.Context(), fromID, toID, req.Amount, req.ServiceType)
	if err != nil {
		logger.Error.Printf("[PayForService] Payment failed from=%d to=%d amount=%.2f: %v", fromID, toID, req.Amount, err)
		respond.Error(w, http.StatusBadRequest, "payment failed", err)
		return
	}

	logger.Info.Printf("[PayForService] Payment success from=%d to=%d amount=%.2f type=%s", fromID, toID, req.Amount, req.ServiceType)
	respond.JSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "payment completed",
	})
}
