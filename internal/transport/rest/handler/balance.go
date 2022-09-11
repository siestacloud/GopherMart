package handler

import (
	"net/http"

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

// * Запрос на списание средств
// * Получение информации о выводе средств
