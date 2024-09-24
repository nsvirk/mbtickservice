package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/nsvirk/mbtickservice/internal/api"
	"github.com/nsvirk/mbtickservice/internal/config"
	"github.com/nsvirk/mbtickservice/internal/logger"
	"github.com/nsvirk/mbtickservice/internal/repository"
	"github.com/nsvirk/mbtickservice/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := repository.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	// Initialize Redis client
	redisClient, err := repository.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize logger
	appLogger := logger.NewAppLogger(db)
	appLogger.Info("App initialized")

	// Initialize ticker service
	tickerService := service.NewTickerService(db, redisClient)
	defer tickerService.Close()
	appLogger.Info("Ticker service initialized")

	// Initialize Echo server
	e := echo.New()
	e.HideBanner = true

	// Initialize API routes
	api.InitRoutes(e, db, tickerService)

	// Start server
	go func() {
		if err := e.Start(":" + cfg.ServerPort); err != nil && err != http.ErrServerClosed {
			appLogger.Error(fmt.Sprintf("Failed to start server: %v", err))
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		appLogger.Error(fmt.Sprintf("Failed to shutdown server: %v", err))
		log.Fatal(err)
	}

	appLogger.Info("Server shut down gracefully")
	log.Println("Server shut down gracefully")
}
