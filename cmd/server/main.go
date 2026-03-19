package main

import (
	"log"
	"net/http"

	"justtype-backend/internal/config"
	"justtype-backend/internal/db"
	internalHttp "justtype-backend/internal/http"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.Connect(cfg.DB_DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Printf("Connected to database successfully")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", internalHttp.HealthHandler(database))

	serverAddr := ":" + cfg.Port
	log.Printf("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
