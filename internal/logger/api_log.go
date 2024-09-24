package logger

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

func RequestLog(db *gorm.DB, userID, botID string, tiCt int) error {
	// Init logger
	logger := NewAppLogger(db, userID, botID)

	// Log
	requestLog := map[string]string{
		"user_id": userID,
		"bot_id":  botID,
		"ti_ct":   fmt.Sprintf("%d", tiCt),
	}
	// Make json
	requestJson, err := json.Marshal(requestLog)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Request: %s", string(requestJson)))
	return nil
}

func ResponseLog(db *gorm.DB, userID, botID, ticksChannel string, subscribedCt int) error {
	// Init logger
	logger := NewAppLogger(db, userID, botID)

	// Log
	responseLog := map[string]string{
		"ticks_channel": ticksChannel,
		"subscribed_ct": fmt.Sprintf("%d", subscribedCt),
	}
	// Make json
	responseJson, err := json.Marshal(responseLog)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Response: %s", string(responseJson)))

	return nil
}
