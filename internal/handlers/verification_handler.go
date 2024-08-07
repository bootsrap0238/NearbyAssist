package handlers

import (
	filehandler "nearbyassist/internal/file"
	"nearbyassist/internal/hash"
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type verificationHandler struct {
	server *server.Server
}

func NewVerificationHandler(s *server.Server) *verificationHandler {
	return &verificationHandler{
		server: s,
	}
}

func (h *verificationHandler) HandleVerifyIdentity(c echo.Context) error {
	name := c.FormValue("name")
	address := c.FormValue("address")
	idType := c.FormValue("idType")
	idNumber := c.FormValue("idNumber")

	req := &models.IdentityVerificationModel{
		Name:     name,
		Address:  address,
		IdType:   idType,
		IdNumber: idNumber,
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse submitted files")
	}

	for _, file := range files {
		handler := filehandler.NewFileHandler(h.server.Encrypt)

		switch file.Filename {
		case "frontId":
			url, err := handler.SavePhoto(file, h.server.Storage.SaveFrontId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save front id")
			}

			if id, err := h.server.DB.NewFrontId(&models.FrontIdModel{Url: url}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save front id to db")
			} else {
				req.FrontId = id
			}
		case "backId":
			url, err := handler.SavePhoto(file, h.server.Storage.SaveBackId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save back id")
			}

			if id, err := h.server.DB.NewBackId(&models.BackIdModel{Url: url}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save back id to db")
			} else {
				req.BackId = id
			}
		case "face":
			url, err := handler.SavePhoto(file, h.server.Storage.SaveFace)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save face")
			}

			if id, err := h.server.DB.NewFace(&models.FaceModel{Url: url}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save face to db")
			} else {
				req.Face = id
			}
		}
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if cipher, err := h.server.Encrypt.EncryptString(req.Name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	} else {
		req.Name = cipher
	}

	if cipher, err := h.server.Encrypt.EncryptString(req.Address); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	} else {
		req.Address = cipher
	}

	if cipher, err := h.server.Encrypt.EncryptString(req.IdNumber); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	} else {
		req.IdNumber = cipher
	}

	verificationId, err := h.server.DB.NewIdentityVerification(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"message":        "Identity verification submitted.",
		"verificationId": verificationId,
	})
}

func (h *verificationHandler) HandleGetAllIdentityVerification(c echo.Context) error {
	requests, err := h.server.DB.FindAllIdentityVerification()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"requests": requests,
	})
}

func (h *verificationHandler) HandleGetIdentityVerification(c echo.Context) error {
	verificationId := c.Param("verificationId")
	id, err := strconv.Atoi(verificationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Verification ID must be a number")
	}

	request, err := h.server.DB.FindIdentityVerificationById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Identity verification not found")
	}

	if decrypted, err := h.server.Encrypt.DecryptString(request.Name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	} else {
		request.Name = decrypted
	}

	if decrypted, err := h.server.Encrypt.DecryptString(request.Address); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	} else {
		request.Address = decrypted
	}

	if decrypted, err := h.server.Encrypt.DecryptString(request.IdNumber); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	} else {
		request.IdNumber = decrypted
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"request": request,
	})
}
