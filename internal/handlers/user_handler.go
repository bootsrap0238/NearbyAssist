package handlers

import (
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	server *server.Server
}

func NewUserHandler(server *server.Server) *userHandler {
	return &userHandler{
		server: server,
	}
}

func (h *userHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "User base route",
	})
}

func (h *userHandler) HandleGetUser(c echo.Context) error {
	userId := c.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	user, err := h.server.DB.FindUserById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"user": user,
	})
}

func (h *userHandler) HandleCount(c echo.Context) error {
	count, err := h.server.DB.CountUser()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"userCount": count,
	})
}
