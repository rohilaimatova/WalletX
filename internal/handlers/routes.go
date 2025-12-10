package handlers

import (
	"WalletX/internal/handlers/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, userHandler *UserHandler, cardHandler *CardHandler, serviceHandler *ServiceHandler, accountHandler *AccountHandler) {

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

	services := api.PathPrefix("").Subrouter()
	services.Use(middleware.CheckUserAuthentication)
	services.HandleFunc("/services", serviceHandler.GetAllServices).Methods("GET")

	api.HandleFunc("/pay", accountHandler.PayForService).Methods("POST")

	secured := api.PathPrefix("").Subrouter() //
	secured.Use(middleware.CheckUserAuthentication)
	cards := secured.PathPrefix("/cards").Subrouter()

	{
		cards.HandleFunc("/create", cardHandler.CreateCard).Methods("POST")
		cards.HandleFunc("/user", cardHandler.GetCardsByUser).Methods("GET")        // ?user_id=
		cards.HandleFunc("/deactivate", cardHandler.DeactivateCard).Methods("POST") // ?card_id=
	}
}
