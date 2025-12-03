package handlers

import (
	"WalletX/internal/handlers/middleware"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, userHandler *UserHandler, accHandler *AccountHandler, txHandler *TransactionHandler) {

	pingHandler := NewHandler() // создаём экземпляр
	r.HandleFunc("/ping", pingHandler.Ping).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()

	// ---- Users (открытые маршруты) ----
	users := api.PathPrefix("/users").Subrouter()
	{
		users.HandleFunc("/signUp", userHandler.SignUp).Methods("POST")
		users.HandleFunc("/set-password", userHandler.SetPassword).Methods("POST")
		users.HandleFunc("/login", userHandler.Login).Methods("POST")
		users.HandleFunc("/verify", userHandler.VerifyIdentity).Methods("POST")
	}

	// ---- Закрытые маршруты с JWT ----
	secured := api.NewRoute().Subrouter()
	secured.Use(middleware.CheckUserAuthentication)

	// ---- Accounts (защищённые) ----
	accounts := secured.PathPrefix("/accounts").Subrouter()
	{
		accounts.HandleFunc("/create", accHandler.CreateAccount).Methods("POST")
		accounts.HandleFunc("/{id}/balance", accHandler.GetBalance).Methods("GET")
	}

	// ---- Transactions (защищённые) ----
	transactions := secured.PathPrefix("/transactions").Subrouter()
	{
		transactions.HandleFunc("/deposit", txHandler.Deposit).Methods("POST")
		transactions.HandleFunc("/withdraw", txHandler.Withdraw).Methods("POST")
		transactions.HandleFunc("/transfer", txHandler.Transfer).Methods("POST")
		transactions.HandleFunc("/{id}/history", txHandler.GetHistory).Methods("GET")
	}
}
