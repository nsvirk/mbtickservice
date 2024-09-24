package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nsvirk/mbtickservice/internal/logger"
	"github.com/nsvirk/mbtickservice/internal/service"
	"github.com/nsvirk/mbtickservice/pkg/response"
	"gorm.io/gorm"
)

// PublishRequest is the request body for the /ticks route
type PublishRequest struct {
	BotID             string   `json:"bot_id"`
	TickerInstruments []string `json:"ticker_instruments"`
}

type PublishResponse struct {
	PublishedChannel string `json:"published_channel,omitempty"`
	SubscribedCount  int    `json:"subscribed_count"`
}

type PublishHandler struct {
	DB            *gorm.DB
	tickerService *service.TickerService
}

func NewPublishHandler(DB *gorm.DB, tickerService *service.TickerService) *PublishHandler {
	return &PublishHandler{DB: DB, tickerService: tickerService}
}

func (h *PublishHandler) PublishTicks(c echo.Context) error {
	db := service.NewDBService(h.DB)

	var req PublishRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "InputException", "Invalid request body")
	}

	if req.BotID == "" || len(req.TickerInstruments) == 0 {
		return response.ErrorResponse(c, http.StatusBadRequest, "InputException", "bot_id and ticker_instruments are required")
	}

	// Parse Authorization header
	auth := c.Request().Header.Get("Authorization")
	parts := strings.SplitN(auth, ":", 2)
	if len(parts) != 2 {
		return response.ErrorResponse(c, http.StatusUnauthorized, "AuthorizationException", "Invalid Authorization header")
	}

	// Get userID and enctoken
	userID, enctoken := parts[0], parts[1]

	// Log Request
	err := logger.RequestLog(h.DB, userID, req.BotID, len(req.TickerInstruments))
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "DatabaseException", "Failed to log request")
	}

	// Get instrument tokens
	instrumentTokenMap, err := db.MakeTickerInstrumentTokenMap(req.TickerInstruments)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "DatabaseException", "Failed to get instrument tokens")
	}

	// Save user connection info
	instrumentCt := len(instrumentTokenMap)
	if err := db.SaveUserConnection(userID, enctoken, instrumentCt); err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "DatabaseException", fmt.Sprintf("Failed to store user connection: %v", err))
	}

	// Save instruments
	if err := db.SaveTickerInstruments(req.BotID, userID, instrumentTokenMap); err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "DatabaseException", "Failed to store instruments")
	}

	// Get tickerInstruments from database
	tickerInstruments, err := db.GetTickerInstruments(req.BotID, userID)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "DatabaseException", fmt.Sprintf("Failed to get instruments: %v", err))
	}

	// Start ticker
	err = h.tickerService.StartTicker(userID, enctoken, req.BotID, tickerInstruments)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "TickerException", fmt.Sprintf("Failed to start ticker: %v", err))
	}

	// Name the channel
	ticksChannel := h.tickerService.GetTicksChannel(userID, req.BotID)

	// Log Response
	err = logger.ResponseLog(h.DB, userID, req.BotID, ticksChannel, instrumentCt)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "DatabaseException", "Failed to log response")
	}

	// Make response
	publishResponse := PublishResponse{
		PublishedChannel: ticksChannel,
		SubscribedCount:  len(tickerInstruments),
	}

	// Send success response
	return response.SuccessResponse(c, publishResponse)
}
