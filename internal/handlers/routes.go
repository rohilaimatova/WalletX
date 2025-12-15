package handlers

import (
	"WalletX/internal/handlers/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(r *mux.Router, userHandler *UserHandler, servicesHandler *ServicesHandler, accountHandler *AccountHandler, userProfileHandler *UserProfileHandler, transferHandler *TransferHandler) {

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
	services.HandleFunc("/services", servicesHandler.GetAllServices).Methods("GET")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.CheckUserAuthentication)
	protected.HandleFunc("/users/profile", userProfileHandler.GetUserProfile).Methods("GET")
	protected.HandleFunc("/users/balance", userProfileHandler.GetUserBalance).Methods("GET")
	protected.HandleFunc("/transfer", transferHandler.Transfer).Methods("POST")
	protected.HandleFunc("/pay", accountHandler.PayForService).Methods("POST")
	protected.HandleFunc("/history", transferHandler.TransactionHistory).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	//http://localhost:8080/swagger/index.html
}
