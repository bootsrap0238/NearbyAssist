package server

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/message/v1"
	"nearbyassist/internal/handlers"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	// middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Custom validator
	e.Validator = &utils.Validator{Validator: validator.New()}

	// File server
	e.Static("/resource", "store/")

	// Routes
	e.GET("/health", health.HealthCheck)
	handlers.RouteHandlerV1(e.Group("/v1"))

	// Goroutine for saving messages
	go message.MessageSavior()

	// Goroutine for sending messages
	go message.MessageForwarder()

	return e
}
