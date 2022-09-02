package handler

import (
	"net/http"
	"strings"

	"github.com/siestacloud/gopherMart/internal/core"

	"github.com/labstack/echo/v4"
)

// @Summary SignUp
// @Tags Auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body core.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/register [post]
func (h *Handler) SignUp() echo.HandlerFunc {

	return func(c echo.Context) error {
		var input core.User

		if err := c.Bind(&input); err != nil {
			return errResponse(c, http.StatusBadRequest, "")
		}

		_, err := h.services.Authorization.CreateUser(input)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				return errResponse(c, http.StatusConflict, "")
			}
			return errResponse(c, http.StatusInternalServerError, "")
		}

		token, err := h.services.Authorization.GenerateToken(input.Login, input.Password)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, "")
		}
		c.Response().Header().Set("Authorization", "Bearer "+token)
		return c.NoContent(http.StatusOK)
	}
}

type signInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignIn
// @Tags Auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/login [post]
func (h *Handler) SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input signInInput
		if err := c.Bind(&input); err != nil {
			return errResponse(c, http.StatusBadRequest, "")
		}

		token, err := h.services.Authorization.GenerateToken(input.Login, input.Password)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				return errResponse(c, http.StatusConflict, "")
			}
			return errResponse(c, http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"token": token,
		})
	}
}
