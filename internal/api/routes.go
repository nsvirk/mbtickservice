package api

import (
	"github.com/labstack/echo/v4"
	"github.com/nsvirk/moneybotstds/internal/api/handlers"
	"github.com/nsvirk/moneybotstds/internal/api/middleware"
	"github.com/nsvirk/moneybotstds/internal/config"
	"github.com/nsvirk/moneybotstds/internal/service"
	"gorm.io/gorm"
)

func InitRoutes(e *echo.Echo, cfg *config.Config, db *gorm.DB, tickerService *service.TickerService) {

	middleware.LoggerMiddleware(e)
	middleware.RecoverMiddleware(e)

	// Create a group for all API routes
	ticks := e.Group("")

	// Create a group for protected routes
	protected := ticks.Group("")
	protected.Use(middleware.AuthMiddleware())

	// Index route
	indexHandler := handlers.NewIndexHandler(cfg)
	e.GET("/", indexHandler.Index)

	// /publish route
	publishHandler := handlers.NewPublishHandler(db, tickerService)
	protected.POST("/publish/start", publishHandler.StartPublishing)
	protected.POST("/publish/stop", publishHandler.StopPublishing)

}
