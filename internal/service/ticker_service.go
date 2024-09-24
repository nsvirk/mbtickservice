package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	kiteticker "github.com/nsvirk/gokiteticker"
	kitemodels "github.com/nsvirk/gokiteticker/models"
	"github.com/nsvirk/mbtickservice/internal/logger"
	"github.com/nsvirk/mbtickservice/internal/models"
	"github.com/nsvirk/mbtickservice/internal/repository"
	"gorm.io/gorm"
)

type TickerService struct {
	db           *gorm.DB
	redisClient  *repository.RedisClient
	tickers      map[string]*TickerInstance
	mu           sync.Mutex
	tickerLogger *logger.TickerLogger
}

type TickerInstance struct {
	Ticker   *kiteticker.Ticker
	TokenMap map[uint32]string
}

type Tick struct {
	Exchange      string
	TradingSymbol string
	PublishedAt   time.Time
	Tick          kitemodels.Tick
}

func NewTickerService(db *gorm.DB, redisClient *repository.RedisClient) *TickerService {
	return &TickerService{
		db:           db,
		redisClient:  redisClient,
		tickers:      make(map[string]*TickerInstance),
		tickerLogger: logger.NewTickerLogger(db),
	}
}

func (s *TickerService) StartTicker(userID, enctoken, botID string, tickerInstruments []models.TickerInstrument) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s:%s", userID, botID)

	// Check if user already has 3 active tickers
	userTickerCount := 0
	for k := range s.tickers {
		if k[:len(userID)] == userID {
			userTickerCount++
		}
	}
	if userTickerCount >= 3 {
		return fmt.Errorf("maximum tickers reached for user %s", userID)
	}

	// Create new Kite ticker instance
	ticker := kiteticker.New(userID, enctoken)

	instance := &TickerInstance{
		Ticker:   ticker,
		TokenMap: make(map[uint32]string),
	}

	// Set up callbacks
	ticker.OnMessage(s.onMessage(userID, botID))
	ticker.OnTick(s.onTick(userID, botID, instance))
	ticker.OnError(s.onError(userID, botID))
	ticker.OnClose(s.onClose(userID, botID))
	ticker.OnConnect(s.onConnect(userID, botID))
	ticker.OnReconnect(s.onReconnect(userID, botID))
	ticker.OnNoReconnect(s.onNoReconnect(userID, botID))

	// Prepare instrument tokens for subscription
	var instTokens []uint32
	for _, inst := range tickerInstruments {
		instTokens = append(instTokens, inst.InstrumentToken)
		instance.TokenMap[inst.InstrumentToken] = fmt.Sprintf("%s:%s", inst.Exchange, inst.Tradingsymbol)
	}

	// Start the connection
	go ticker.Serve()

	// Wait for the connection to be established
	connectionEstablished := make(chan bool)
	ticker.OnConnect(func() {
		s.onConnect(userID, botID)()
		connectionEstablished <- true
	})

	select {
	case <-connectionEstablished:
		fmt.Println("Connection established")
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timed out waiting for ticker connection")
	}

	// Subscribe to instruments
	if err := ticker.Subscribe(instTokens); err != nil {
		return fmt.Errorf("subscription error: %w", err)
	}

	// Set subscription mode
	if err := ticker.SetMode(kiteticker.ModeFull, instTokens); err != nil {
		return fmt.Errorf("setMode error: %w", err)
	}

	// Store ticker instance
	s.tickers[key] = instance

	// Log the event
	s.logTickerEvent(userID, botID, "INFO", "StartTicker", "Ticker started successfully")

	return nil
}

func (s *TickerService) StopTicker(userID, botID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s:%s", userID, botID)

	instance, exists := s.tickers[key]
	if !exists {
		return fmt.Errorf("ticker not found for user %s and bot %s", userID, botID)
	}

	// Unsubscribe from all tokens
	var tokens []uint32
	for token := range instance.TokenMap {
		tokens = append(tokens, token)
	}

	if len(tokens) > 0 {
		if err := instance.Ticker.Unsubscribe(tokens); err != nil {
			s.logTickerEvent(userID, botID, "ERROR", "StopTicker", fmt.Sprintf("Failed to unsubscribe: %v", err))
			// Continue with stopping even if unsubscribe fails
		}
	}

	// Stop the ticker
	instance.Ticker.Stop()

	// Close the connection
	if err := instance.Ticker.Close(); err != nil {
		s.logTickerEvent(userID, botID, "ERROR", "StopTicker", fmt.Sprintf("Failed to close connection: %v", err))
	}

	// Remove the ticker instance from the map
	delete(s.tickers, key)

	// Log the event
	s.logTickerEvent(userID, botID, "INFO", "StopTicker", "Ticker stopped successfully")

	return nil
}

func (s *TickerService) onTick(userID, botID string, instance *TickerInstance) func(tick kitemodels.Tick) {
	return func(tick kitemodels.Tick) {
		// Get exchange and tradingsymbol
		instrument, ok := instance.TokenMap[tick.InstrumentToken]
		if !ok {
			s.logTickerEvent(userID, botID, "ERROR", "onTick", fmt.Sprintf("Unknown instrument token: %d", tick.InstrumentToken))
			return
		}

		parts := strings.Split(instrument, ":")
		if len(parts) != 2 {
			s.logTickerEvent(userID, botID, "ERROR", "onTick", fmt.Sprintf("Invalid instrument: %s", instrument))
			return
		}

		exchange := parts[0]
		tradingSymbol := parts[1]

		// For new Tick
		newTick := Tick{
			Exchange:      exchange,
			TradingSymbol: tradingSymbol,
			PublishedAt:   time.Now(),
			Tick:          tick,
		}

		tickJSON, err := json.Marshal(newTick)
		if err != nil {
			s.logTickerEvent(userID, botID, "ERROR", "onTick", fmt.Sprintf("Failed to marshal tick: %v", err))
			return
		}

		// Publish tick to Redis channel
		channelName := fmt.Sprintf("CH:TICKS:%s:%s", userID, botID)
		if err := s.redisClient.PublishTicks(channelName, tickJSON); err != nil {
			s.logTickerEvent(userID, botID, "ERROR", "PublishTicks", fmt.Sprintf("Failed to publish tick: %v", err))
		}
	}
}

func (s *TickerService) onError(userID, botID string) func(err error) {
	return func(err error) {
		s.logTickerEvent(userID, botID, "ERROR", "onError", err.Error())
	}
}

func (s *TickerService) onClose(userID, botID string) func(code int, reason string) {
	return func(code int, reason string) {
		s.logTickerEvent(userID, botID, "INFO", "onClose", fmt.Sprintf("Connection closed: code=%d, reason=%s", code, reason))
	}
}

func (s *TickerService) onConnect(userID, botID string) func() {
	return func() {
		s.logTickerEvent(userID, botID, "INFO", "onConnect", "Connected to Kite ticker")
	}
}

func (s *TickerService) onReconnect(userID, botID string) func(attempt int, delay time.Duration) {
	return func(attempt int, delay time.Duration) {
		s.logTickerEvent(userID, botID, "INFO", "onReconnect", fmt.Sprintf("Reconnected to Kite ticker after %d attempts, delay: %v", attempt, delay))
	}
}
func (s *TickerService) onMessage(userID, botID string) func(messageType int, message []byte) {
	return func(messageType int, message []byte) {
		if messageType == 1 {
			s.logTickerEvent(userID, botID, "INFO", "onMessage", fmt.Sprintf("Received message: type=%d, message=%s", messageType, string(message)))
		}
	}
}

func (s *TickerService) onNoReconnect(userID, botID string) func(attempt int) {
	return func(attempt int) {
		s.logTickerEvent(userID, botID, "INFO", "onNoReconnect", fmt.Sprintf("No reconnect after %d attempts", attempt))
	}
}

func (s *TickerService) logTickerEvent(userID, botID, level, eventType, message string) {
	s.tickerLogger.Log(userID, botID, level, eventType, message)
}

func (s *TickerService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, instance := range s.tickers {
		instance.Ticker.Close()
	}
}

func (s *TickerService) GetTicksChannel(userID, botID string) string {
	return fmt.Sprintf("CH:TICKS:%s:%s", userID, botID)
}
