package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// Cargar .env antes de leer las variables
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: .env file not found, using system environment variables")
	}
}

// GetEnvironment returns the current environment
func GetEnvironment() string {
	return os.Getenv("ENVIRONMENT")
}

var (
	JWTSecret       = []byte("your_secret_key_change_in_production")
	TokenExpiration = 1 * time.Hour
	ServerPort      = ":8080"
)
