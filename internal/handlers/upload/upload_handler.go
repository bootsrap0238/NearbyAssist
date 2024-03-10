package upload

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/photo/v1"

	"github.com/labstack/echo/v4"
)

func UploadHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck)
	r.POST("/upload", photo.UploadImage)
}
