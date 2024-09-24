package api

import (
	"github.com/labstack/echo/v4"
	"github.com/nsvirk/mbtickservice/internal/api/handlers"
	"github.com/nsvirk/mbtickservice/internal/api/middleware"
	"github.com/nsvirk/mbtickservice/internal/service"
	"gorm.io/gorm"
)

func InitRoutes(e *echo.Echo, db *gorm.DB, tickerService *service.TickerService) {

	middleware.LoggerMiddleware(e)
	middleware.RecoverMiddleware(e)

	// Create a group for all API routes
	ticks := e.Group("/ticks")

	// Create a group for protected routes
	protected := ticks.Group("")
	protected.Use(middleware.AuthMiddleware())

	// /publish route
	publishHandler := handlers.NewPublishHandler(db, tickerService)
	protected.POST("/publish", publishHandler.PublishTicks)

}
