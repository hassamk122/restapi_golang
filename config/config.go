package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	DatabaseUrl string
	Environment string
	LogLevel    string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("Error loading env file %v", err)
	}
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseUrl: getEnv("DATABASE_URL", "postgres"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, defaulValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaulValue
}
