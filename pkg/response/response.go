package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response represents the standard API response structure
type Response struct {
	Status    string      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	ErrorType string      `json:"error_type,omitempty"`
	Message   string      `json:"message,omitempty"`
}

type TicksResponse struct {
	Status string `json:"status"`
	Data   struct {
		UserID       string `json:"user_id"`
		BotID        string `json:"bot_id"`
		TicksChannel string `json:"ticks_channel"`
		SubscribedCt int    `json:"subscribed_ct"`
	} `json:"data"`
}

// SuccessResponse sends a successful JSON response
func SuccessResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Status: "ok",
		Data:   data,
	})
}

// ErrorResponse sends an error JSON response
func ErrorResponse(c echo.Context, httpStatus int, errorType, message string) error {
	return c.JSON(httpStatus, Response{
		Status:    "error",
		ErrorType: errorType,
		Message:   message,
	})
}
