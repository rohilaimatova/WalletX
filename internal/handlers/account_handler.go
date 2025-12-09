package handlers

/*
import (
	"WalletX/internal/service"
	"WalletX/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	AccountService     *service.AccountService
	TransactionService *service.TransactionService
}

func NewAccountHandler(acc *service.AccountService, tx *service.TransactionService) *AccountHandler {
	return &AccountHandler{
		AccountService:     acc,
		TransactionService: tx,
	}
}
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	acc, err := h.AccountService.CreateAccount(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(acc)
}

func (h *AccountHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	balance, err := h.AccountService.GetBalance(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(models.BalanceResponse{
		AccountID: id,
		Balance:   balance,
	})
}
*/
