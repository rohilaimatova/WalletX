package main

import (
	"WalletX/config"
	"WalletX/internal/db"
	"WalletX/internal/handlers"
	"WalletX/internal/repository"
	"WalletX/internal/service"
	"WalletX/pkg/logger"
	redisPkg "WalletX/pkg/redis"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	if err := config.ReadSettings(); err != nil {
		logger.Error.Fatalf("Ошибка чтения настроек: %s", err)
	}

	if err := logger.Init(); err != nil {
		logger.Error.Fatalf("Error initializing logger: %v", err)
	}

	rdb := redisPkg.InitRedis()
	if rdb == nil {
		logger.Error.Fatalf("Ошибка подключения Redis")
	}

	if err := db.ConnectDB(); err != nil {
		logger.Error.Fatalf("Error connecting to DB: %v", err)
	}
	defer db.CloseDB()

	conn := db.GetDBConnection()

	userRepo := repository.NewPostgresUserRepo(conn)
	//accRepo := repository.NewPostgresAccountRepo(conn)
	//txRepo := repository.NewPostgresTransactionRepo(conn)
	cardRepo := repository.NewPostgresCardRepo(conn)
	accountRepo := repository.NewAccountRepository(conn)
	serviceRepo := repository.NewServiceRepo(conn)

	userService := service.NewUserService(userRepo)
	//accService := service.NewAccountService(accRepo)
	//txService := service.NewTransactionService(accRepo, txRepo)
	cardService := service.NewCardService(cardRepo) // ← сервис карт
	accountService := service.NewAccountService(accountRepo)
	serviceService := service.NewServiceService(serviceRepo) // создаём сервис аккаунта

	userHandler := handlers.NewUserHandler(userService, accountService, rdb)
	//accHandler := handlers.NewAccountHandler(accService, txService)
	//txHandler := handlers.NewTransactionHandler(accService, txService)
	cardHandler := handlers.NewCardHandler(cardService)
	serviceHandler := handlers.NewServiceHandler(serviceService)

	r := mux.NewRouter()
	handlers.RegisterRoutes(r, userHandler, cardHandler, serviceHandler) // ← передаем cardHandler

	logger.Info.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Error.Fatalf("Server error: %v", err)
	}
}
