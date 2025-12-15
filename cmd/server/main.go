package main

import (
	"WalletX/config"
	"WalletX/internal/db"
	"WalletX/internal/handlers"
	"WalletX/internal/handlers/transaction"
	"WalletX/internal/repository"
	"WalletX/internal/service"
	"WalletX/pkg/logger"
	redisPkg "WalletX/pkg/redis"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	if err := config.ReadSettings(); err != nil {
		logger.Error.Fatalf("Failed to read configuration settings: %s", err)
	}

	if err := logger.Init(); err != nil {
		logger.Error.Fatalf("Error initializing logger: %v", err)
	}

	rdb := redisPkg.InitRedis()
	if rdb == nil {
		logger.Error.Fatalf("Failed to connect to Redis")
	}

	if err := db.ConnectDB(); err != nil {
		logger.Error.Fatalf("Error connecting to DB: %v", err)
	}
	defer db.CloseDB()

	conn := db.GetDBConnection()

	userRepo := repository.NewPostgresUserRepo(conn)
	accountRepo := repository.NewAccountRepository(conn)
	servicesRepo := repository.NewServicesRepo(conn)
	transactionRepo := repository.NewTransactionRepository(conn)
	profileRepo := repository.NewUserProfileRepository(conn)
	transactionManager := transaction.NewTransactionManager(conn)

	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	servicesService := service.NewServicesService(servicesRepo)
	userProfileService := service.NewUserProfileService(profileRepo)
	transferService := service.NewTransferService(accountRepo, transactionRepo, transactionManager)
	paymentService := service.NewPaymentService(accountRepo, transactionRepo, servicesRepo, transactionManager)

	userHandler := handlers.NewUserHandler(userService, accountService, rdb)
	servicesHandler := handlers.NewServicesHandler(servicesService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	userProfileHandler := handlers.NewUserProfileHandler(userProfileService)
	transferHandler := handlers.NewTransferHandler(transferService)

	r := mux.NewRouter()
	handlers.RegisterRoutes(r, userHandler, servicesHandler, paymentHandler, userProfileHandler, transferHandler)

	logger.Info.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Error.Fatalf("Server error: %v", err)
	}
}
