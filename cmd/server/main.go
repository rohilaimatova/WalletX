package main

import (
	"WalletX/config"
	"WalletX/internal/db"
	"WalletX/internal/handlers"
	"WalletX/internal/repository"
	"WalletX/internal/service"
	"WalletX/pkg/logger"
	redisPkg "WalletX/pkg/redis" // <-- ВАЖНО!
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	// ---- CONFIG ----
	if err := config.ReadSettings(); err != nil {
		log.Fatalf("Ошибка чтения настроек: %s", err)
	}

	// ---- LOGGER ----
	if err := logger.Init(); err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}

	// ---- INIT REDIS ----
	rdb := redisPkg.InitRedis()

	// ---- DB CONNECTION ----
	if err := db.ConnectDB(); err != nil {
		logger.Error.Printf("Error connecting to DB: %v", err)
		return
	}
	defer db.CloseDB()

	conn := db.GetDBConnection()

	// ---- REPOSITORIES ----
	userRepo := repository.NewPostgresUserRepo(conn)
	accRepo := repository.NewPostgresAccountRepo(conn)
	txRepo := repository.NewPostgresTransactionRepo(conn)

	// ---- SERVICES ----
	userService := service.NewUserService(userRepo)
	accService := service.NewAccountService(accRepo)
	txService := service.NewTransactionService(accRepo, txRepo)

	// ---- HANDLERS ----
	userHandler := handlers.NewUserHandler(userService, rdb) // <-- передаём Redis!
	accHandler := handlers.NewAccountHandler(accService, txService)
	txHandler := handlers.NewTransactionHandler(accService, txService)

	// ---- ROUTER ----
	r := mux.NewRouter()
	handlers.RegisterRoutes(r, userHandler, accHandler, txHandler)

	// ---- START SERVER ----
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
