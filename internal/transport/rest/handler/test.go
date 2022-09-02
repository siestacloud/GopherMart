package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// @Summary Test
// @Tags Test
// @Security ApiKeyAuth
// @Description test handler
// @ID testingID
// @Accept  json
// @Produce  json
// @Success 200 {string} string "OK"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/test/ [get]
func (h *Handler) Test() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Info("Ok")
		// h.services.Authorization.Test()
		return c.HTML(http.StatusOK, "OK")
	}
}
