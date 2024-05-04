package handlers

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type serviceHandler struct {
	server *server.Server
}

func NewServiceHandler(server *server.Server) *serviceHandler {
	return &serviceHandler{
		server: server,
	}
}

func (h *serviceHandler) HandleGetServices(c echo.Context) error {
	services, err := h.server.DB.FindAllService()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}

func (h *serviceHandler) HandleRegisterService(c echo.Context) error {
	req := &request.NewService{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	authHeader := c.Request().Header.Get("Authorization")
	if userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		req.VendorId = userId
	}

	models.ConstructLocationFromLatLong(&req.GeoSpatialModel)

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	insertId, err := h.server.DB.RegisterService(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"serviceId": insertId,
	})
}

func (h *serviceHandler) HandleUpdateService(c echo.Context) error {
	serviceId := c.Param("serviceId")
	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	req := &request.UpdateService{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
	}

	authHeader := c.Request().Header.Get("Authorization")
	if userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		req.VendorId = userId
	}

	models.ConstructLocationFromLatLong(&req.GeoSpatialModel)

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	} else {
		req.Id = id
	}

	if err := h.server.DB.UpdateService(req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not process update request")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Update service",
	})
}

func (h *serviceHandler) HandleDeleteService(c echo.Context) error {
	serviceId := c.Param("serviceId")
	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	if err := h.server.DB.DeleteService(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *serviceHandler) HandleSearchService(c echo.Context) error {
	params, err := utils.GetSearchParams(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	services, err := h.server.DB.GeoSpatialSearch(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}

func (h *serviceHandler) HandleGetDetails(c echo.Context) error {
	serviceId := c.Param("serviceId")
	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	service, err := h.server.DB.FindServiceById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	vendor, err := h.server.DB.FindVendorByService(service.ServiceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: retrieve review count

	images, err := h.server.DB.FindAllPhotosByServiceId(service.ServiceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"serviceInfo":   service,
		"vendorInfo":    vendor,
		"serviceImages": images,
	})
}

func (h *serviceHandler) HandleGetByVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "owner ID must be a number")
	}

	services, err := h.server.DB.FindServiceByVendor(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"services": services,
	})
}

// takes origin as QueryString ex: origin=lat,long
func (h *serviceHandler) HandleFindRoute(c echo.Context) error {
	serviceId := c.Param("serviceId")
	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	service, err := h.server.DB.FindServiceById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not find service")
	}

	origin, err := parseOrigin(c.QueryParam("origin"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid origin coordinates")
	}

	distination := models.NewLocationWithData(service.Latitude, service.Longitude)

	polyline, err := h.server.RouteEngine.FindRoute(origin, distination)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not find routes at the moment")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"polyline": polyline,
	})
}

func parseOrigin(query string) (*models.Location, error) {
	coords := strings.Split(query, ",")
	lat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, err
	}

	long, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, err
	}

	origin := models.NewLocationWithData(lat, long)

	return origin, nil
}
