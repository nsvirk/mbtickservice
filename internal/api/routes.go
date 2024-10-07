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
	api := e.Group("")

	// Index route
	indexHandler := handlers.NewIndexHandler(cfg)
	e.GET("/", indexHandler.Index)

	// /publish route
	publishHandler := handlers.NewPublishHandler(db, tickerService)
	publishGroup := api.Group("/publish")
	publishGroup.Use(middleware.AuthMiddleware())
	publishGroup.POST("/start", publishHandler.StartPublishing)
	publishGroup.POST("/stop", publishHandler.StopPublishing)

}
