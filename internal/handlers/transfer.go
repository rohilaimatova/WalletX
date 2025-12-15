package handlers

import (
	"WalletX/internal/handlers/middleware"
	"WalletX/internal/service"
	"WalletX/models"
	"WalletX/pkg/logger"
	"WalletX/respond"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type TransferHandler struct {
	TransferService *service.TransferService
}

func NewTransferHandler(ts *service.TransferService) *TransferHandler {
	return &TransferHandler{
		TransferService: ts,
	}
}

func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req models.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn.Printf("[TransferHandler] Invalid request body: %v", err)
		respond.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	userIDRaw := r.Context().Value(middleware.UserIDCtx)
	if userIDRaw == nil {
		logger.Warn.Println("[TransferHandler] User not authenticated")
		respond.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}
	fromUserID := userIDRaw.(int)

	fromAcc, err := h.TransferService.AccountRepo.GetByUserID(r.Context(), fromUserID)
	if err != nil {
		logger.Warn.Printf("[TransferHandler] Sender account not found: userID=%d", fromUserID)
		respond.Error(w, http.StatusBadRequest, "sender account not found", err)
		return
	}

	toAcc, err := h.TransferService.AccountRepo.GetByPhone(r.Context(), req.ToPhone)
	if err != nil {
		logger.Warn.Printf("[TransferHandler] Recipient account not found: phone=%s", req.ToPhone)
		respond.Error(w, http.StatusNotFound, "recipient not found", err)
		return
	}

	if err := h.TransferService.Transfer(r.Context(), fromAcc.ID, toAcc.ID, req.Amount); err != nil {
		if err.Error() == "cannot transfer to your own account" {
			logger.Warn.Printf("[TransferHandler] Attempt to transfer to self: fromID=%d", fromAcc.ID)
			respond.Error(w, http.StatusBadRequest, "cannot transfer to your own account", errors.New("cannot transfer to your own account"))
			return
		}
		logger.Error.Printf("[TransferHandler] Transfer failed: %v", err)
		respond.Error(w, http.StatusBadRequest, "transfer failed", err)
		return
	}

	logger.Info.Printf("[TransferHandler] Transfer completed: fromAccountID=%d, toAccountID=%d, amount=%.2f",
		fromAcc.ID, toAcc.ID, req.Amount)
	respond.JSON(w, http.StatusOK, map[string]string{"status": "success"})
}
func (h *TransferHandler) TransactionHistory(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("[TransferHandler] TransactionHistory called")

	query := r.URL.Query()
	startStr := query.Get("start")
	endStr := query.Get("end")

	layout := "2006-01-02"
	start, err := time.Parse(layout, startStr)
	if err != nil || startStr == "" {
		start = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		logger.Info.Printf("[TransferHandler] No valid start date provided, using %s", start)
	}

	end, err := time.Parse(layout, endStr)
	if err != nil || endStr == "" {
		end = time.Now()
		logger.Info.Printf("[TransferHandler] No valid end date provided, using %s", end)
	} else {
		end = end.Add(24*time.Hour - time.Nanosecond)
	}

	userIDRaw := r.Context().Value(middleware.UserIDCtx)
	if userIDRaw == nil {
		logger.Warn.Println("[TransferHandler] User not authenticated")
		respond.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}
	userID := userIDRaw.(int)
	logger.Info.Printf("[TransferHandler] UserID from token: %d", userID)

	account, err := h.TransferService.AccountRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		logger.Warn.Printf("[TransferHandler] Account not found for userID=%d: %v", userID, err)
		respond.Error(w, http.StatusBadRequest, "account not found", err)
		return
	}
	logger.Info.Printf("[TransferHandler] Found account: %+v", account)

	transactions, err := h.TransferService.AccountRepo.GetTransactions(r.Context(), account.ID, start, end)
	if err != nil {
		logger.Error.Printf("[TransferHandler] Failed to get transactions for accountID=%d: %v", account.ID, err)
		respond.Error(w, http.StatusInternalServerError, "failed to get transactions", err)
		return
	}
	logger.Info.Printf("[TransferHandler] Retrieved %d transactions for accountID=%d", len(transactions), account.ID)

	respond.JSON(w, http.StatusOK, transactions)
}
