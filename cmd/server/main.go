package main

import (
	"queueit/internal/api"
	"queueit/internal/config"
	"queueit/internal/db"
	"queueit/pkg/logger"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		logger.Fatal(err)
	}

	if err := db.InitDB(); err != nil {
		logger.Fatal(err)
	}

	router := api.NewRouter()
	if err := router.StartServer(); err != nil {
		logger.Fatal(err)
	}
}
