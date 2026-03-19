package http

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type SyncUserRequest struct {
	Email string  `json:"email"`
	Name  *string `json:"name"`
	Image *string `json:"image"`
}

func HealthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := db.Ping(); err != nil {
			log.Printf("Health check failed (DB ping): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "error"})
			return
		}

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

func SyncUserHandler(db *sql.DB, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 1. Secret Check
		if r.Header.Get("X-Internal-Secret") != secret {
			log.Printf("Unauthorized sync attempt from %s", r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}

		// 2. Body Parsing
		var req SyncUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid_json"})
			return
		}

		if req.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "email_required"})
			return
		}

		// 3. Exact Upsert Logic
		query := `
			INSERT INTO users (email, name, image)
			VALUES ($1, $2, $3)
			ON CONFLICT (email)
			DO UPDATE SET name = EXCLUDED.name, image = EXCLUDED.image
		`
		_, err := db.Exec(query, req.Email, req.Name, req.Image)
		if err != nil {
			log.Printf("Failed to upsert user %s: %v", req.Email, err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "db_error"})
			return
		}

		log.Printf("Successfully synced user: %s", req.Email)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
