package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	UserServiceURL  string
	OrderServiceURL string
}

func Load() Config {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	return Config{
		Port:            getEnv("API_GATEWAY_PORT", "8080"),
		UserServiceURL:  getEnv("USER_SERVICE_URL", "localhost:50051"),
		OrderServiceURL: getEnv("ORDER_SERVICE_URL", "localhost:50052"),
	}
}

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
