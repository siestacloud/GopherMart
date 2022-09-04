package handler

import (
	"io/ioutil"
	"net/http"
	"strconv"
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
		userID, err := getUserId(c)
		if err != nil {

			pkg.ErrPrint("transport", http.StatusInternalServerError, err)
			return errResponse(c, http.StatusInternalServerError, err.Error()) // в контексте нет id пользователя
		}

		var order core.Order
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			pkg.ErrPrint("transport", http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "bad request")
		}

		order.ID, err = strconv.Atoi(string(body))
		if err != nil {
			pkg.ErrPrint("transport", http.StatusUnprocessableEntity, err)
			return errResponse(c, http.StatusUnprocessableEntity, "order format failure")
		}

		if err := c.Validate(order); err != nil {
			pkg.ErrPrint("transport", http.StatusBadRequest, err)
			return errResponse(c, http.StatusBadRequest, "validate failure")
		}

		if err := h.services.Order.Create(userID, order); err != nil {
			if strings.Contains(err.Error(), "lune") {
				pkg.ErrPrint("transport", http.StatusUnprocessableEntity, err)
				return errResponse(c, http.StatusUnprocessableEntity, "order format failure")
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

		pkg.InfoPrint("transport", "accepted", order.ID)
		return c.NoContent(http.StatusAccepted)
	}
}

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
		return c.NoContent(http.StatusOK)
	}
}
