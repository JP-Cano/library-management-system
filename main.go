package main

import (
	"fmt"
	"github.com/juanPabloCano/library-management-system/src/api"
	"github.com/juanPabloCano/library-management-system/src/config"
	"github.com/juanPabloCano/library-management-system/src/database"
	"log"
	"os"
)

func main() {
	dbURL := config.DatabaseUrlBuilder()
	psqlStorage := config.NewPsqlStorage(dbURL)
	defer psqlStorage.Close()

	db := database.NewStorage(psqlStorage.Conn)
	_, err := psqlStorage.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	server := api.New(fmt.Sprintf(":%s", serverPort), db)
	if err := server.Serve(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
