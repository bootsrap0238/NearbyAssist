package server

import (
	"nearbyassist/internal/db"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) HandleVersionOneRoutes(r *echo.Group) {

	r.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	r.POST("/register", s.HandleRegister)

	r.POST("/login", s.HandleLogin)

	r.GET("/locations", func(c echo.Context) error {
		locations, err := db.GetLocations()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"locations": locations,
		})
	})
}
