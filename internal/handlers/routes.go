package handlers

import (
	"WalletX/internal/handlers/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, userHandler *UserHandler, cardHandler *CardHandler, serviceHandler *ServiceHandler) {

	pingHandler := NewHandler()
	r.HandleFunc("/ping", pingHandler.Ping).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()

	users := api.PathPrefix("/users").Subrouter()
	{
		users.HandleFunc("/signUp", userHandler.SignUp).Methods("POST")
		users.HandleFunc("/set-password", userHandler.SetPassword).Methods("POST")
		users.HandleFunc("/login", userHandler.Login).Methods("POST")
		users.HandleFunc("/verify", userHandler.VerifyIdentity).Methods("POST")
	}

	services := api.PathPrefix("/services").Subrouter()
	services.Use(middleware.CheckUserAuthentication) // Добавляем проверку токена
	services.HandleFunc("", serviceHandler.GetAllServices).Methods("GET")

	secured := api.PathPrefix("").Subrouter() // упрощённый путь вместо NewRoute
	secured.Use(middleware.CheckUserAuthentication)
	cards := secured.PathPrefix("/cards").Subrouter()

	{
		cards.HandleFunc("/create", cardHandler.CreateCard).Methods("POST")
		cards.HandleFunc("/user", cardHandler.GetCardsByUser).Methods("GET")        // ?user_id=
		cards.HandleFunc("/deactivate", cardHandler.DeactivateCard).Methods("POST") // ?card_id=
	}
	// ---- Accounts (protected) ----
	/*accounts := secured.PathPrefix("/accounts").Subrouter()
	{
		accounts.HandleFunc("/create", accHandler.CreateAccount).Methods("POST")
		accounts.HandleFunc("/{id}/balance", accHandler.GetBalance).Methods("GET")
	}

	// ---- Transactions (protected) ----
	transactions := secured.PathPrefix("/transactions").Subrouter()
	{
		transactions.HandleFunc("/deposit", txHandler.Deposit).Methods("POST")
		transactions.HandleFunc("/withdraw", txHandler.Withdraw).Methods("POST")
		transactions.HandleFunc("/transfer", txHandler.Transfer).Methods("POST")
		transactions.HandleFunc("/{id}/history", txHandler.GetHistory).Methods("GET")
	}

	*/
}
