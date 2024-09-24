package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL      string
	PostgresSchema   string
	PostgresLogLevel string
	RedisAddr        string
	ServerPort       string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		PostgresURL:      getEnv("TS_PG_DSN", ""),
		PostgresSchema:   getEnv("TS_PG_SCHEMA", ""),
		PostgresLogLevel: getEnv("TS_PG_LOG_LEVEL", "error"),
		RedisAddr:        getEnv("TS_REDIS_ADDR", ""),
		ServerPort:       getEnv("TS_SERVER_PORT", ""),
	}

	if config.PostgresURL == "" {
		return nil, fmt.Errorf("TS_PG_DSN is required")
	}

	if config.PostgresSchema == "" {
		return nil, fmt.Errorf("TS_PG_SCHEMA is required")
	}

	if config.RedisAddr == "" {
		return nil, fmt.Errorf("TS_REDIS_ADDR is required")
	}

	if config.ServerPort == "" {
		return nil, fmt.Errorf("TS_SERVER_PORT is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
