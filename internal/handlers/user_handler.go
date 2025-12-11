package handlers

import (
	"WalletX/internal/handlers/middleware"
	"WalletX/internal/service"
	"WalletX/pkg/logger"
	"WalletX/respond"
	"net/http"
)

type UserProfileHandler struct {
	Service *service.UserProfileService
}

func NewUserProfileHandler(s *service.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		Service: s,
	}
}

func (h *UserProfileHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.UserIDCtx)
	if userIDRaw == nil {
		respond.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		respond.Error(w, http.StatusInternalServerError, "invalid userID type", nil)
		return
	}

	user, err := h.Service.GetProfileByID(userID)
	if err != nil {
		respond.HandleError(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, user)
	logger.Info.Printf("User profile returned for ID %d", userID)
}

func (h *UserProfileHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.UserIDCtx)
	if userIDRaw == nil {
		respond.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		respond.Error(w, http.StatusInternalServerError, "invalid userID type", nil)
		return
	}

	balance, err := h.Service.GetBalanceByUserID(userID)
	if err != nil {
		respond.HandleError(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, balance)
	logger.Info.Printf("Balance returned for user ID %d", userID)
}
