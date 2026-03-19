package main

import (
	"log"
	"net/http"

	"justtype-backend/internal/config"
	"justtype-backend/internal/db"
	internalHttp "justtype-backend/internal/http"
)

var allowedOrigins = []string{
	"http://localhost:3000",
	"http://localhost:8080",
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Only add CORS headers if Origin is in allowlist
		for _, o := range allowedOrigins {
			if origin == o {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Internal-Secret")
				break
			}
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	mux.HandleFunc("/internal/sync-user", internalHttp.SyncUserHandler(database, cfg.InternalAPISecret))

	// Wrap server with CORS middleware
	handler := corsMiddleware(mux)

	serverAddr := ":" + cfg.Port
	log.Printf("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
