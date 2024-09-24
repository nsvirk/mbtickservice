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
