package logger

import (
	"log"
	"time"

	"github.com/nsvirk/mbtickservice/internal/models"
	"gorm.io/gorm"
)

// TickerLogger is the main struct for the ticker logger
type TickerLogger struct {
	db *gorm.DB
}

// NewTickerLogger creates a new ticker	logger
func NewTickerLogger(db *gorm.DB) *TickerLogger {
	return &TickerLogger{
		db: db,
	}
}

// Log a TICKER message
func (l *TickerLogger) Log(userID, botID, level, eventType, message string) {
	tickerLog := models.TickerLog{
		Timestamp: time.Now(),
		UserID:    userID,
		BotID:     botID,
		Level:     level,
		EventType: eventType,
		Message:   message,
	}
	err := l.db.Table(models.TickerLogsTable).Create(&tickerLog).Error
	if err != nil {
		log.Println("Error logging TICKER message:", err)
	}
}

// TickerInstanceLogger is the main struct for the ticker instance logger
type TickerInstanceLogger struct {
	db *gorm.DB
}

// NewTickerInstanceLogger creates a new ticker instance logger
func NewTickerInstanceLogger(db *gorm.DB) *TickerInstanceLogger {
	return &TickerInstanceLogger{db: db}
}

// Log a TICKER_INSTANCE message
func (l *TickerInstanceLogger) Log(userID, botID, instance string) {
	tickerInstanceLog := models.TickerInstanceLog{
		UserID:   userID,
		BotID:    botID,
		Instance: instance,
	}
	err := l.db.Table(models.TickerInstanceLogsTable).Create(&tickerInstanceLog).Error
	if err != nil {
		log.Println("Error logging TICKER INSTANCE:", err)
	}
}
