package main

import (
	"WalletX/config"
	"WalletX/internal/db"
	"WalletX/internal/handlers"
	"WalletX/internal/repository"
	"WalletX/internal/service"
	"WalletX/pkg/logger"
	redisPkg "WalletX/pkg/redis"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	// ---- CONFIG ----
	if err := config.ReadSettings(); err != nil {
		logger.Error.Fatalf("Ошибка чтения настроек: %s", err)
	}

	// ---- LOGGER ----
	if err := logger.Init(); err != nil {
		logger.Error.Fatalf("Error initializing logger: %v", err)
	}

	// ---- INIT REDIS ----
	rdb := redisPkg.InitRedis()
	if rdb == nil {
		logger.Error.Fatalf("Ошибка подключения Redis")
	}

	// ---- DB CONNECTION ----
	if err := db.ConnectDB(); err != nil {
		logger.Error.Fatalf("Error connecting to DB: %v", err)
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
	userHandler := handlers.NewUserHandler(userService, rdb)
	accHandler := handlers.NewAccountHandler(accService, txService)
	txHandler := handlers.NewTransactionHandler(accService, txService)

	// ---- ROUTER ----
	r := mux.NewRouter()
	handlers.RegisterRoutes(r, userHandler, accHandler, txHandler)

	// ---- START SERVER ----
	logger.Info.Println("Server running on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Error.Fatalf("Server error: %v", err)
	}
}
