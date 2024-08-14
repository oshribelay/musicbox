package main

import (
	"github.com/joho/godotenv"
	"github.com/oshribelay/musicbox/internal/logger"
	"github.com/oshribelay/musicbox/internal/server"
)

func main() {
	logger.InfoLogger.Println("Started program!")

	if err := godotenv.Load(".env"); err != nil {
		logger.ErrorLogger.Println((err.Error()))
	}

	if err := server.Start(); err != nil {
		logger.ErrorLogger.Println((err.Error()))
	}

	

}