package handler

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/pkg"
)

// @Summary CreateOrder
// @Security ApiKeyAuth
// @Tags Order
// @Description create and validate client order
// @Accept  text/plain
// @Produce  text/plain
// @Param input body integer true "new title and description for item"
// @Success 200,202 {int} int "no content"
// @Failure 400,401,422 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/orders [post]
func (h *Handler) CreateOrder() echo.HandlerFunc {

	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "new request", "/api/user/orders")
		userID, err := getUserId(c)
		if err != nil {
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, err.Error()) // в контексте нет id пользователя
		}
		// * парсинг
		var order core.Order
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			pkg.ErrPrint("transport", http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "bad request")
		}
		order.Number = string(body)
		// * валидация номера заказа
		if err := c.Validate(order); err != nil {
			pkg.ErrPrint("transport", http.StatusUnprocessableEntity, err)
			return errResponse(c, http.StatusUnprocessableEntity, "order format failure")
		}
		// * получаю информацию о расчете начислений баллов лояльности (внешнее api)
		if err := h.services.Accrual.GetOrderAccrual(&order); err != nil {
			return errResponse(c, http.StatusBadRequest, err.Error())
		}
		// * проверяю заказ по алг луна и добавляю в бд (связывая с клиентом)
		if err := h.services.Order.Create(userID, order); err != nil {
			if strings.Contains(err.Error(), "lune") {
				pkg.ErrPrint("transport", http.StatusUnprocessableEntity, err)
				return errResponse(c, http.StatusUnprocessableEntity, err.Error())
			}
			if strings.Contains(err.Error(), "user already have order") {
				pkg.InfoPrint("transport", "ok", err)
				return c.NoContent(http.StatusOK)
			}
			if strings.Contains(err.Error(), "another user order") {
				pkg.ErrPrint("transport", http.StatusConflict, err)
				return errResponse(c, http.StatusConflict, err.Error())
			}
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		// * получаю информацию о расчете начислений баллов лояльности (внешнее api)
		if err := h.services.Balance.UpdateCurrent(userID, &order); err != nil {
			return errResponse(c, http.StatusBadRequest, err.Error())
		}

		pkg.InfoPrint("transport", "accepted", order.Number)
		return c.NoContent(http.StatusAccepted)
	}
}

// * Получение списка полученных номеров заказов
// @Summary GetOrder
// @Security ApiKeyAuth
// @Tags Order
// @Description create and validate client order
// @ID get_order
// @Accept  text/plain
// @Produce  text/plain
// @Success 200,202 {int} int "no content"
// @Failure 400,401,422 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/user/orders [get]
func (h *Handler) GetOrders() echo.HandlerFunc {
	return func(c echo.Context) error {
		pkg.InfoPrint("transport", "new request", "/api/user/orders")

		userID, err := getUserId(c)
		if err != nil {
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, err.Error()) // в контексте нет id пользователя
		}

		orderList, err := h.services.GetListOrders(userID)
		if err != nil {
			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, "internal server error")
		}

		respList := []core.Order{}
		for i, v := range orderList {
			// * получаю информацию о расчете начислений баллов лояльности (внешнее api)
			if err := h.services.Accrual.GetOrderAccrual(&orderList[i]); err != nil {
				pkg.ErrPrint("transport", http.StatusInternalServerError, err, v)
				// return errResponse(c, http.StatusBadRequest, err.Error())
			}
			if orderList[i].Status == "" {
				orderList[i].Status = "PROCESSING"
			}
		}

		c.Request().Header.Set("Content-Type", "application/json")

		respList = append(respList, orderList[len(orderList)-1])
		if len(respList) == 0 {
			pkg.ErrPrint("transport", http.StatusNoContent, "no data to answer")
			return errResponse(c, http.StatusNoContent, "")
		}
		pkg.InfoPrint("transport", "OK", respList)
		return c.JSON(http.StatusOK, respList)

	}
}
