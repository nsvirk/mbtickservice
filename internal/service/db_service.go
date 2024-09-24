package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/nsvirk/mbtickservice/internal/models"
	"github.com/nsvirk/mbtickservice/internal/repository"
	"gorm.io/gorm"
)

type DBService struct {
	repo *repository.Repository
}

func NewDBService(db *gorm.DB) *DBService {
	return &DBService{repo: repository.NewRepository(db)}
}

func (s *DBService) SaveUserConnection(userID, enctoken string, instrumentCt int) error {
	// Create a new user
	user := models.User{
		UserID:        userID,
		Enctoken:      enctoken[0:8] + "******",
		InstrumentsCt: instrumentCt,
		ConnectedAt:   time.Now(),
	}

	return s.repo.UpsertUser(&user)
}

func (s *DBService) MakeTickerInstrumentTokenMap(tickerInstruments []string) (map[string]uint32, error) {
	instrumentTokenMap := make(map[string]uint32)
	for _, instrument := range tickerInstruments {
		parts := strings.Split(instrument, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid instrument format: %s", instrument)
		}
		exchange, tradingsymbol := parts[0], parts[1]

		instrumentToken, err := s.repo.GetInstrumentToken(exchange, tradingsymbol)
		if err != nil {
			return nil, fmt.Errorf("error querying instrument token: %w", err)
		}

		instrumentTokenMap[instrument] = instrumentToken
	}

	return instrumentTokenMap, nil
}

func (s *DBService) SetTickerInstruments(botID, userID string, instrumentTokenMap map[string]uint32) error {

	// Delete all ticker instruments
	err := s.repo.DeleteTickerInstruments(botID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete ticker instruments: %w", err)
	}

	// Insert new ticker instruments
	tickerInstruments := make([]models.TickerInstrument, 0, len(instrumentTokenMap))
	for instrument, token := range instrumentTokenMap {
		parts := strings.Split(instrument, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid instrument format: %s", instrument)
		}
		exchange, tradingsymbol := parts[0], parts[1]

		now := time.Now()
		tickerInstrument := models.TickerInstrument{
			BotID:           botID,
			UserID:          userID,
			Exchange:        exchange,
			Tradingsymbol:   tradingsymbol,
			InstrumentToken: token,
			UpdatedAt:       now,
		}

		tickerInstruments = append(tickerInstruments, tickerInstrument)
	}

	return s.repo.InsertTickerInstruments(tickerInstruments)
}

// Get TickerInstruments from the database
func (s *DBService) GetTickerInstruments(botID, userID string) ([]models.TickerInstrument, error) {
	return s.repo.GetTickerInstruments(botID, userID)
}
