package handlers

import (
	"errors"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
	server *server.Server
}

func NewAuthHandler(server *server.Server) *authHandler {
	return &authHandler{server}
}

func (h *authHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Auth route is up and running!",
	})
}

func (h *authHandler) HandleAdminLogin(c echo.Context) error {
	req := &request.AdminLogin{}
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	admin, err := h.server.DB.FindAdminByUsername(req.Username)
	if admin == nil || err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	accessToken, err := h.server.Auth.GenerateAdminAccessToken(admin)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshToken, err := h.server.Auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	session := models.NewSessionModel(refreshToken)
	if _, err := h.server.DB.NewSession(session); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"role":         admin.Role,
		"adminId":      admin.Id,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *authHandler) HandleLogin(c echo.Context) error {
	req := &request.UserLogin{}
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userModel := models.NewUserModelWithData(req.Name, req.Email, req.Image)

	// Check if user exists
	user, err := h.server.DB.FindUserByEmail(req.Email)
	if err != nil {
		// Register user if not found
		userId, err := h.server.DB.NewUser(req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		userModel.Id = userId
	} else {
		userModel.Id = user.Id
	}

	accessToken, err := h.server.Auth.GenerateUserAccessToken(userModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshToken, err := h.server.Auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	session := models.NewSessionModel(refreshToken)
	if _, err := h.server.DB.NewSession(session); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"userId":       userModel.Id,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *authHandler) HandleLogout(c echo.Context) error {
	req := &request.Logout{}
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	session, err := h.server.DB.FindActiveSessionByToken(req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Session not found")
	}

	if err := h.server.DB.LogoutSession(session.Id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := h.server.DB.BlacklistToken(session.Token); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Logout successful",
	})
}

func (h *authHandler) HandleTokenRefresh(c echo.Context) error {
	req := &request.RefreshToken{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err := h.server.DB.FindSessionByToken(req.Token); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid token")
	}

	if blacklist, _ := h.server.DB.FindBlacklistedToken(req.Token); blacklist != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Token blacklisted")
	}

	authHeader := c.Request().Header.Get("Authorization")
	userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	jwtClaims, err := h.server.Auth.GetClaims(authHeader[len("Bearer "):])
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse JWT")
	}

	var newAccessToken string // If there's an err, user is not an admin
	if _, err := utils.GetRoleFromClaims(jwtClaims); err != nil {
		user, err := h.server.DB.FindUserById(userId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if token, err := h.server.Auth.GenerateUserAccessToken(user); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			newAccessToken = token
		}
	} else {
		admin, err := h.server.DB.FindAdminById(userId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if token, err := h.server.Auth.GenerateAdminAccessToken(admin); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			newAccessToken = token
		}
	}

	// As defense against my dumb self, check first if newAccessToken is an empty string
	if newAccessToken == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("Some error occurred while generating new access token"))
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"accessToken": newAccessToken,
	})
}
