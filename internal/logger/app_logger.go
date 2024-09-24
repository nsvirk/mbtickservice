package logger

import (
	"log"
	"time"

	"github.com/nsvirk/mbtickservice/internal/models"
	"gorm.io/gorm"
)

// AppLogger is the main struct for the logger
type AppLogger struct {
	db     *gorm.DB
	userID string
	botID  string
}

// NewLogger creates a new logger
func NewAppLogger(db *gorm.DB, userID, botID string) *AppLogger {
	return &AppLogger{db: db, userID: userID, botID: botID}
}

// Log logs a message
func (l *AppLogger) Log(level, message string) error {
	now := time.Now()
	log := models.Log{
		Timestamp: now,
		UserID:    l.userID,
		BotID:     l.botID,
		Level:     level,
		Message:   message,
	}
	return l.db.Table(models.LogsTable).Create(&log).Error
}

// Log a INFO message
func (l *AppLogger) Info(message string) {
	err := l.Log("INFO", message)
	if err != nil {
		log.Println("Error logging INFO message:", err)
	}
}

// Log a WARN message
func (l *AppLogger) Warn(message string) {
	err := l.Log("WARN", message)
	if err != nil {
		log.Println("Error logging WARN message:", err)
	}
}

// Log a ERROR message
func (l *AppLogger) Error(message string) {
	err := l.Log("ERROR", message)
	if err != nil {
		log.Println("Error logging ERROR message:", err)
	}
}

// Log a FATAL message
func (l *AppLogger) Fatal(message string) {
	err := l.Log("FATAL", message)
	if err != nil {
		log.Println("Error logging FATAL message:", err)
	}
}

// Log a DEBUG message
func (l *AppLogger) Debug(message string) {
	err := l.Log("DEBUG", message)
	if err != nil {
		log.Println("Error logging DEBUG message:", err)
	}
}
