package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() Config {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	cfg := Config{
		GRPCPort:   getEnv("USER_SERVICE_GRPC_PORT", "50051"),
		DBHost:     getEnv("USER_DB_HOST", "localhost"),
		DBPort:     getEnv("USER_DB_PORT", "5432"),
		DBUser:     getEnv("USER_DB_USER", "postgres"),
		DBPassword: getEnv("USER_DB_PASSWORD", "postgres"),
		DBName:     getEnv("USER_DB_NAME", "user_db"),
	}

	return cfg
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
	)
}

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
