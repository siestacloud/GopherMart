package handler

import (
	"net/http"
	"strings"

	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/pkg"

	"github.com/labstack/echo/v4"
)

// 	* `POST /api/user/register` 					— регистрация пользователя;
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
			pkg.ErrPrint("transport", http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "bind body failure")
		}
		if err := c.Validate(input); err != nil {
			pkg.ErrPrint("transport", http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "validate failure")
		}

		// * авторизация
		userID, err := h.services.Authorization.CreateUser(input)
		if err != nil {
			if strings.Contains(err.Error(), "login busy") {
				pkg.ErrPrint("transport", http.StatusConflict, err)
				return errResponse(c, http.StatusConflict, err.Error())
			}

			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		// * аутентификация
		token, err := h.services.Authorization.GenerateToken(input.Login, input.Password)
		if err != nil {

			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")

		}
		// * создание баланса для нового клиента
		if err := h.services.Balance.Create(userID); err != nil {
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		c.Response().Header().Set("Authorization", "Bearer "+token)
		return c.NoContent(http.StatusOK)
	}
}

type signInInput struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// 	* `POST /api/user/login` 						— аутентификация пользователя;
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
			pkg.ErrPrint("transport", http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "bind body failure")
		}
		if err := c.Validate(input); err != nil {
			pkg.ErrPrint("transport", http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "validate failure")
		}
		token, err := h.services.Authorization.GenerateToken(input.Login, input.Password)
		if err != nil {
			if strings.Contains(err.Error(), "invalid username/password pair") {
				pkg.ErrPrint("transport", http.StatusBadRequest, err)
				return errResponse(c, http.StatusUnauthorized, err.Error())
			}

			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}
		c.Response().Header().Set("Authorization", "Bearer "+token)
		return c.NoContent(http.StatusOK)
	}
}
