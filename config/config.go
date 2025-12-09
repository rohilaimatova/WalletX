package config

import (
	"WalletX/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var AppSettings models.Config

func ReadSettings() error {

	fmt.Println("Loading .env file")

	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found, using system environment variables")
	}

	wd, err := os.Getwd()
	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't get working directory: %s", err.Error()))
	}

	configPath := filepath.Join(wd, "config", "config.json")
	fmt.Println("Reading settings file:", configPath)

	configFile, err := os.Open(configPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't open config file: %s", err.Error()))
	}
	defer func(configFile *os.File) {
		err = configFile.Close()
		if err != nil {
			log.Fatal("Couldn't close config file: ", err.Error())
		}
	}(configFile)

	if err = json.NewDecoder(configFile).Decode(&AppSettings); err != nil {
		return errors.New(fmt.Sprintf("Couldn't decode json config: %s", err.Error()))
	}

	return nil
}
