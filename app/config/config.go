package config

import (
	"github.com/gofiber/fiber/v2/log"
	"os"

	_ "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	JWTSecret     []byte
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		// Jika file .env tidak ditemukan, lanjutkan tanpa error
		log.Fatalf("Error loading .env file")
	}
	return &Config{
		ServerAddress: ":8080",
		JWTSecret:     []byte(os.Getenv("JWT_SECRET")),
	}
}
