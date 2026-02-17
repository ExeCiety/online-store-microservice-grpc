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

	return Config{
		GRPCPort:   getEnv("ORDER_SERVICE_GRPC_PORT", "50052"),
		DBHost:     getEnv("ORDER_DB_HOST", "localhost"),
		DBPort:     getEnv("ORDER_DB_PORT", "5433"),
		DBUser:     getEnv("ORDER_DB_USER", "postgres"),
		DBPassword: getEnv("ORDER_DB_PASSWORD", "postgres"),
		DBName:     getEnv("ORDER_DB_NAME", "order_db"),
	}
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
