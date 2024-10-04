package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/nsvirk/moneybotstds/internal/config"
	"github.com/nsvirk/moneybotstds/pkg/response"
)

// IndexHandler is the handler for the /publish routes
type IndexHandler struct {
	cfg *config.Config
}

// NewIndexHandler creates a new IndexHandler
func NewIndexHandler(cfg *config.Config) *IndexHandler {
	return &IndexHandler{cfg: cfg}
}

func (h *IndexHandler) Index(c echo.Context) error {
	return response.SuccessResponse(c, map[string]string{
		"app_name":    h.cfg.AppName,
		"app_version": h.cfg.AppVersion,
	})
}
