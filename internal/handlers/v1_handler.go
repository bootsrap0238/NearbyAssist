package handlers

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/handlers/auth"
	"nearbyassist/internal/handlers/service"
	"nearbyassist/internal/handlers/user"

	"github.com/labstack/echo/v4"
)

func RouteHandlerV1(r *echo.Group) {
	r.GET("/health", health.HealthCheck)

	authGroup := r.Group("/auth")
	auth.AuthHandler(authGroup)

	serviceGroup := r.Group("/services")
	service.ServiceHandler(serviceGroup)

	userGroup := r.Group("/users")
	user.UserHandler(userGroup)
}