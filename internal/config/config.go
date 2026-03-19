package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN            string
	Port              string
	InternalAPISecret string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // Ignore error if .env is missing

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DB_DSN:            os.Getenv("DB_DSN"),
		Port:              port,
		InternalAPISecret: os.Getenv("INTERNAL_API_SECRET"),
	}
}
