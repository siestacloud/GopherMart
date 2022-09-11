package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/siestacloud/gopherMart/pkg"
)

// * Получение текущего баланса пользователя
// @Summary GetBalance
// @Security ApiKeyAuth
// @Tags Balance
// @Description check client balance
// @ID get_balance
// @Accept  text/plain
// @Produce  text/plain
// @Success 200 {int} int "no content"
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/balance [get]
func (h *Handler) GetBalance() echo.HandlerFunc {
	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "new request", "/api/user/balance")

		userID, err := getUserId(c)
		if err != nil {
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, err.Error()) // в контексте нет id пользователя
		}

		userBalance, err := h.services.Balance.Get(userID)
		if err != nil {
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		c.Request().Header.Set("Content-Type", "application/json")

		pkg.InfoPrint("transport", "OK", userBalance)
		return c.JSON(http.StatusOK, userBalance)

	}
}

type withdrawInput struct {
	Order string  `json:"order" validate:"required"`
	Sum   float64 `json:"sum" validate:"required"`
}

// * Запрос на списание средств
// @Summary WithdrawBalance
// @Security ApiKeyAuth
// @Tags Withdraw
// @Description Withdraw user balance
// @Accept  json
// @Produce  json
// @Param input WithdrawInput  true "some description"
// @Success 200{int} int "no content"
// @Failure 401,402,422 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/balance/withdraw [post]
func (h *Handler) WithdrawBalance() echo.HandlerFunc {
	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "new request", "/api/user/balance/withdrawals")
		userID, err := getUserId(c)
		if err != nil {
			pkg.ErrPrint("transport", http.StatusUnauthorized, err)
			return errResponse(c, http.StatusUnauthorized, err.Error()) // в контексте нет id пользователя
		}

		var input withdrawInput
		if err := c.Bind(&input); err != nil {
			pkg.ErrPrint("transport", http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "bind body failure")
		}
		fmt.Println("ORDER  ", input.Order)

		if err := pkg.Valid(input.Order); err != nil {
			pkg.ErrPrint("transport", http.StatusUnprocessableEntity, err)
			return errResponse(c, http.StatusUnprocessableEntity, "invalid order number")
		}

		if err := h.services.Withdrawal(userID, input.Sum); err != nil {
			if strings.Contains(err.Error(), "there are not enough points on the balance") {
				pkg.ErrPrint("transport", http.StatusPaymentRequired, err)
				return errResponse(c, http.StatusPaymentRequired, err.Error())
			}
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "unable withdraw balls from balance")
		}

		pkg.InfoPrint("transport", "accepted", userID)
		return c.NoContent(http.StatusOK)
	}
}

// * Получение информации о выводе средств
// @Summary WithdrawalsBalance
// @Security ApiKeyAuth
// @Tags WithdrawalsBalance
// @Description check client WithdrawalsBalance
// @ID get_WithdrawalsBalance
// @Accept  text/plain
// @Produce  text/plain
// @Success 200,204 {int} int "no content"
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/balance/withdrawals [get]
func (h *Handler) WithdrawalsBalance() echo.HandlerFunc {
	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "new request", "/api/user/balance/withdrawals")

		userID, err := getUserId(c)
		if err != nil {
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, err.Error()) // в контексте нет id пользователя
		}

		userBalance, err := h.services.Balance.Get(userID)
		if err != nil {
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		c.Request().Header.Set("Content-Type", "application/json")

		pkg.InfoPrint("transport", "OK", userBalance)
		return c.NoContent(http.StatusOK)

	}
}
