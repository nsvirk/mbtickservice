package repository

import (
	"fmt"

	"github.com/nsvirk/mbtickservice/internal/config"
	"github.com/nsvirk/mbtickservice/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection
func InitDB(cfg *config.Config) (*gorm.DB, error) {

	// Set log level
	var logLevel logger.LogLevel
	switch cfg.PostgresLogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(cfg.PostgresURL), &gorm.Config{Logger: logger.Default.LogMode(logLevel)})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create schema
	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", cfg.PostgresSchema)
	tx := db.Exec(sql)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to create schema: %w", tx.Error)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.TickerInstrument{}, &models.Log{}, &models.TickerLog{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Check if db is init
	if db == nil {
		return nil, fmt.Errorf("failed to init database: %w", err)
	}

	return db, nil
}
