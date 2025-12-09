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

type TransactionHandler struct {
	AccountService     *service.AccountService
	TransactionService *service.TransactionService
}

func NewTransactionHandler(acc *service.AccountService, tx *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		AccountService:     acc,
		TransactionService: tx,
	}
}

func (h *TransactionHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	var req models.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tx, err := h.TransactionService.Deposit(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(tx)
}

func (h *TransactionHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req models.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tx, err := h.TransactionService.Withdraw(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(tx)
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req models.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tx, err := h.TransactionService.Transfer(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(tx)
}

func (h *TransactionHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	history, err := h.TransactionService.GetHistory(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(history)
}
*/
