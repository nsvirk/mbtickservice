package models

import (
	"log"
	"time"

	"github.com/nsvirk/moneybotstds/internal/config"
)

var (
	SchemaName             = getSchemaName()
	UsersTable             = SchemaName + "." + "users"
	TickerInstrumentsTable = SchemaName + "." + "ticker_instruments"
	LogsTable              = SchemaName + "." + "logs"
	TickerLogsTable        = SchemaName + "." + "ticker_logs"
	InstrumentsTable       = "api.instruments"
)

func getSchemaName() string {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return cfg.PostgresSchema
}

// User represents the tickserver.users table
type User struct {
	ID            uint32 `gorm:"primaryKey"`
	UserID        string `gorm:"uniqueIndex"`
	Enctoken      string
	InstrumentsCt int
	ConnectedAt   time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return UsersTable
}

// TickerInstrument represents the ticker instruments table
type TickerInstrument struct {
	ID              uint32 `gorm:"primaryKey"`
	UserID          string `gorm:"index;idx_user_bot_token,priority:1"`
	BotID           string `gorm:"index;idx_user_bot_token,priority:2"`
	InstrumentToken uint32 `gorm:"index;idx_user_bot_token,priority:3"`
	Exchange        string
	Tradingsymbol   string
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (TickerInstrument) TableName() string {
	return TickerInstrumentsTable
}

// Log represents the logs table
type Log struct {
	ID        uint32 `gorm:"primaryKey"`
	Timestamp time.Time
	Level     string
	Message   string
}

func (Log) TableName() string {
	return LogsTable
}

// TickerLog represents the ticker logs table
type TickerLog struct {
	ID        uint32 `gorm:"primaryKey"`
	Timestamp time.Time
	UserID    string
	BotID     string
	Level     string
	EventType string
	Message   string
}

func (TickerLog) TableName() string {
	return TickerLogsTable
}
