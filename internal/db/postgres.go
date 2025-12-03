package db

import (
	"WalletX/config"
	"WalletX/pkg/logger"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var db *sql.DB

func ConnectDB() error {
	cfg := config.AppSettings.PostgresParams

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		os.Getenv("DB_PASSWORD"),
		cfg.Database,
	)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		logger.Error.Printf("[db] ConnectDB():error during conect to postgres:%s ", err.Error())
		return err
	}
	return nil

}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
func GetDBConnection() *sql.DB {
	return db
}
