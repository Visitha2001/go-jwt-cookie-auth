package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// EnvConfig loads a variable from .env file or environment
func EnvConfig(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using system env")
	}
	return os.Getenv(key)
}
