package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nsvirk/mbtickservice/internal/service"
	"github.com/nsvirk/mbtickservice/pkg/response"
	"gorm.io/gorm"
)

// StartPublishRequest is the request body for the /publish/start route
type StartPublishRequest struct {
	BotID             string   `json:"bot_id"`
	TickerInstruments []string `json:"ticker_instruments"`
}

// StopPublishRequest is the request body for the /publish/stop route
type StopPublishRequest struct {
	BotID string `json:"bot_id"`
}

// StartPublishResponse is the response body for the /publish/start route
type StartPublishResponse struct {
	PublishedChannel string `json:"published_channel,omitempty"`
	SubscribedCount  int    `json:"subscribed_count"`
}

// StopPublishResponse is the response body for the /publish/stop route
type StopPublishResponse struct {
	Message string `json:"message"`
}

// PublishHandler is the handler for the /publish routes
type PublishHandler struct {
	DB            *gorm.DB
	tickerService *service.TickerService
}

// NewPublishHandler creates a new PublishHandler
func NewPublishHandler(DB *gorm.DB, tickerService *service.TickerService) *PublishHandler {
	return &PublishHandler{DB: DB, tickerService: tickerService}
}

// StartPublishing starts the publishing of ticks
func (h *PublishHandler) StartPublishing(c echo.Context) error {
	db := service.NewDBService(h.DB)

	var req StartPublishRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "InputException", "Invalid request body")
	}

	if req.BotID == "" || len(req.TickerInstruments) == 0 {
		return response.ErrorResponse(c, http.StatusBadRequest, "InputException", "`bot_id` and `ticker_instruments` are required")
	}

	// Parse Authorization header
	auth := c.Request().Header.Get("Authorization")
	parts := strings.SplitN(auth, ":", 2)
	if len(parts) != 2 {
		return response.ErrorResponse(c, http.StatusUnauthorized, "AuthorizationException", "Invalid Authorization header")
	}

	// Get userID and enctoken
	userID, enctoken := parts[0], parts[1]

	// Get instrument tokens from the database
	instrumentTokenMap, err := db.MakeTickerInstrumentTokenMap(req.TickerInstruments)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "DatabaseException", "Failed to get instrument tokens")
	}

	// Save user info to the database
	instrumentCt := len(instrumentTokenMap)
	if err := db.SaveUserConnection(userID, enctoken, instrumentCt); err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "DatabaseException", fmt.Sprintf("Failed to store user connection: %v", err))
	}

	// Set ticker instruments in the database
	if err := db.SetTickerInstruments(req.BotID, userID, instrumentTokenMap); err != nil {
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

	// Get the channel name
	ticksChannel := h.tickerService.GetTicksChannel(userID, req.BotID)

	// Make response
	startPublishResponse := StartPublishResponse{
		PublishedChannel: ticksChannel,
		SubscribedCount:  len(tickerInstruments),
	}

	// Send success response
	return response.SuccessResponse(c, startPublishResponse)
}

// StopPublishing stops the publishing of ticks
func (h *PublishHandler) StopPublishing(c echo.Context) error {

	var req StopPublishRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "InputException", "Invalid request body")
	}

	if req.BotID == "" {
		return response.ErrorResponse(c, http.StatusBadRequest, "InputException", "`bot_id` is required")
	}

	// Parse Authorization header
	auth := c.Request().Header.Get("Authorization")
	parts := strings.SplitN(auth, ":", 2)
	if len(parts) != 2 {
		return response.ErrorResponse(c, http.StatusUnauthorized, "AuthorizationException", "Invalid Authorization header")
	}

	// Get userID
	userID := parts[0]

	// Stop ticker
	err := h.tickerService.StopTicker(userID, req.BotID)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "TickerException", fmt.Sprintf("Failed to stop ticker: %v", err))
	}

	// Make response
	stopPublishResponse := StopPublishResponse{
		Message: "Publishing stopped successfully",
	}

	// Send success response
	return response.SuccessResponse(c, stopPublishResponse)
}
