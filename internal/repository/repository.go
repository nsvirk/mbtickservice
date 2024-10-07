package repository

import (
	"errors"
	"fmt"

	"github.com/nsvirk/moneybotstds/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetUser by userID
func (r *Repository) GetUser(userID string) (*models.User, error) {
	var user models.User
	err := r.db.Table(models.UsersTable).Where("user_id = ?", userID).First(&user).Error
	return &user, err
}

// UpsertUser - insert or update a user
func (r *Repository) UpsertUser(user *models.User) error {
	result := r.db.Table(models.UsersTable).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"enctoken", "instruments_ct", "connected_at"}),
	}).Create(&user)

	return result.Error
}

// DeleteTickerInstruments - delete all ticker instruments
func (r *Repository) DeleteTickerInstruments(botID, userID string) error {
	return r.db.Table(models.TickerInstrumentsTable).Where("bot_id = ? AND user_id = ?", botID, userID).Delete(&models.TickerInstrument{}).Error
}

// InsertTickerInstruments - insert multiple ticker instruments
func (r *Repository) InsertTickerInstruments(tickerInstruments []models.TickerInstrument) error {
	return r.db.
		Table(models.TickerInstrumentsTable).
		Create(&tickerInstruments).
		Error
}

// UpsertTickerInstrument - insert or update a ticker instrument
func (r *Repository) UpsertTickerInstrument(tickerInstrument models.TickerInstrument) error {
	err := r.db.Table(models.TickerInstrumentsTable).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "bot_id"}, {Name: "user_id"}, {Name: "instrument_token"}},
		DoUpdates: clause.AssignmentColumns([]string{"exchange", "tradingsymbol", "updated_at"}),
	}).Create(&tickerInstrument).Error

	if err != nil {
		return fmt.Errorf("failed to upsert instrument: %w", err)
	}

	return nil
}

// GetInstrumentToken- get the instrument token
func (r *Repository) GetInstrumentToken(exchange, tradingsymbol string) (uint32, error) {
	var instrumentToken uint32
	err := r.db.Table(models.InstrumentsTable).
		Select("instrument_token").
		Where("exchange = ? AND tradingsymbol = ?", exchange, tradingsymbol).
		Scan(&instrumentToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("instrument not found: %w", err)
		}
		return 0, fmt.Errorf("error querying instrument token: %w", err)
	}
	return instrumentToken, nil
}

// GetTickerInstruments - get the instruments from the API
func (r *Repository) GetTickerInstruments(botID, userID string) ([]models.TickerInstrument, error) {
	var tickerInstruments []models.TickerInstrument
	err := r.db.
		Table(models.TickerInstrumentsTable).
		Where("user_id = ? AND bot_id = ?", userID, botID).
		Find(&tickerInstruments).Error
	return tickerInstruments, err
}
