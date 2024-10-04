package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppName          string
	AppVersion       string
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
		AppName:          getEnv("MB_TDS_APP_NAME", "Moneybots Tick Data Service"),
		AppVersion:       getEnv("MB_TDS_APP_VERSION", "1.0.0"),
		PostgresURL:      getEnv("MB_TDS_PG_DSN", ""),
		PostgresSchema:   getEnv("MB_TDS_PG_SCHEMA", ""),
		PostgresLogLevel: getEnv("MB_TDS_PG_LOG_LEVEL", "error"),
		RedisHost:        getEnv("MB_TDS_REDIS_HOST", ""),
		RedisPort:        getEnv("MB_TDS_REDIS_PORT", ""),
		RedisPassword:    getEnv("MB_TDS_REDIS_PASSWORD", ""),
		ServerPort:       getEnv("MB_TDS_SERVER_PORT", ""),
	}

	if config.PostgresURL == "" {
		return nil, fmt.Errorf("MB_TDS_PG_DSN is required")
	}

	if config.PostgresSchema == "" {
		return nil, fmt.Errorf("MB_TDS_PG_SCHEMA is required")
	}

	if config.RedisHost == "" {
		return nil, fmt.Errorf("MB_TDS_REDIS_HOST is required")
	}

	if config.RedisPort == "" {
		return nil, fmt.Errorf("MB_TDS_REDIS_PORT is required")
	}

	if config.RedisPassword == "" {
		return nil, fmt.Errorf("MB_TDS_REDIS_PASSWORD is required")
	}

	if config.ServerPort == "" {
		return nil, fmt.Errorf("MB_TDS_SERVER_PORT is required")
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
