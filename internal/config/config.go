package config

import (
	"fmt"
	"os"
)

type Config struct {
	PostgresURL      string
	PostgresSchema   string
	PostgresLogLevel string
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	ServerPort       string
}

func Load() (*Config, error) {

	config := &Config{
		PostgresURL:      getEnv("TS_PG_DSN", ""),
		PostgresSchema:   getEnv("TS_PG_SCHEMA", ""),
		PostgresLogLevel: getEnv("TS_PG_LOG_LEVEL", "error"),
		RedisHost:        getEnv("TS_REDIS_HOST", ""),
		RedisPort:        getEnv("TS_REDIS_PORT", ""),
		RedisPassword:    getEnv("TS_REDIS_PASSWORD", ""),
		ServerPort:       getEnv("TS_SERVER_PORT", ""),
	}

	if config.PostgresURL == "" {
		return nil, fmt.Errorf("TS_PG_DSN is required")
	}

	if config.PostgresSchema == "" {
		return nil, fmt.Errorf("TS_PG_SCHEMA is required")
	}

	if config.RedisHost == "" {
		return nil, fmt.Errorf("TS_REDIS_HOST is required")
	}

	if config.RedisPort == "" {
		return nil, fmt.Errorf("TS_REDIS_PORT is required")
	}

	if config.RedisPassword == "" {
		return nil, fmt.Errorf("TS_REDIS_PASSWORD is required")
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
