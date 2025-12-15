package logger

import (
	"WalletX/config"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
	Debug *log.Logger
)

func Init() error {
	logParams := config.AppSettings.LogParams

	logDir := logParams.LogDirectory
	absPath, err := filepath.Abs(logDir)
	if err != nil {
		return fmt.Errorf("cannot resolve log directory: %w", err)
	}

	err = os.MkdirAll(absPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create log directory: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("cannot stat log directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", absPath)
	}

	newLumberjack := func(filename string) *lumberjack.Logger {
		return &lumberjack.Logger{
			Filename:   filepath.Join(absPath, filename), // файл внутри папки logs
			MaxSize:    logParams.MaxSizeMegabytes,
			MaxBackups: logParams.MaxBackups,
			MaxAge:     logParams.MaxAgeDays,
			Compress:   logParams.Compress,
			LocalTime:  logParams.LocalTime,
		}
	}

	Info = log.New(io.MultiWriter(os.Stdout, newLumberjack(logParams.LogInfo)), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stdout, newLumberjack(logParams.LogError)), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(io.MultiWriter(os.Stdout, newLumberjack(logParams.LogWarn)), "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(io.MultiWriter(os.Stdout, newLumberjack(logParams.LogDebug)), "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
