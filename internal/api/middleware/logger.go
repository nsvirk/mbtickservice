// Package middleware provides the middleware for the Echo instance
package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// LoggerMiddleware configures and adds logger middleware to the Echo instance
func LoggerMiddleware(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}: ip=${remote_ip}, req=${method}, uri=${uri}, status=${status}\n",
	}))

}

// RecoverMiddleware configures and adds recover middleware to the Echo instance
func RecoverMiddleware(e *echo.Echo) {
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))
}
