package handlers

import (
	"WalletX/internal/service"
	"WalletX/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type CardHandler struct {
	Service *service.CardService
}

func NewCardHandler(s *service.CardService) *CardHandler {
	return &CardHandler{Service: s}
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var cardReq models.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&cardReq); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(cardReq.CardNumber) != 16 {
		http.Error(w, "card number must be 16 digits", http.StatusBadRequest)
		return
	}

	if len(cardReq.CVV) != 3 {
		http.Error(w, "cvv must be 3 digits", http.StatusBadRequest)
		return
	}

	masked := cardReq.CardNumber[:4] + "****" + cardReq.CardNumber[12:]

	created, err := h.Service.CreateCard(
		cardReq.UserID,
		cardReq.AccountID,
		cardReq.CardType,
		masked, // маска вместо полного номера
		cardReq.ExpiryDate,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *CardHandler) GetCardsByUser(w http.ResponseWriter, r *http.Request) {
	// userID берем, например, из query params: ?user_id=1
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userID := 0
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	cards, err := h.Service.GetCardsByUser(userID)
	if err != nil {
		http.Error(w, "failed to get cards", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cards)
}

func (h *CardHandler) DeactivateCard(w http.ResponseWriter, r *http.Request) {
	// cardID берем, например, из query params: ?card_id=1
	cardIDStr := r.URL.Query().Get("card_id")
	if cardIDStr == "" {
		http.Error(w, "card_id is required", http.StatusBadRequest)
		return
	}

	cardID := 0
	_, err := fmt.Sscanf(cardIDStr, "%d", &cardID)
	if err != nil {
		http.Error(w, "invalid card_id", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeactivateCard(cardID); err != nil {
		http.Error(w, "failed to deactivate card", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "card deactivated"})
}
