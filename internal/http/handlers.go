package http

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func HealthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := db.Ping(); err != nil {
			log.Printf("Health check failed (DB ping): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "error"})
			return
		}

		// Additional lightweight SELECT 1 check
		var result int
		if err := db.QueryRow("SELECT 1").Scan(&result); err != nil {
			log.Printf("Health check failed (SELECT 1): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "error"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
