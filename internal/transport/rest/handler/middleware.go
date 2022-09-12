package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/siestacloud/gopherMart/pkg"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		header := c.Request().Header.Get(authorizationHeader)
		if header == "" {
			return errResponse(c, http.StatusUnauthorized, "empty auth header")
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return errResponse(c, http.StatusUnauthorized, "invalid auth header")
		}

		if len(headerParts[1]) == 0 {
			return errResponse(c, http.StatusUnauthorized, "token is empty")
		}

		userID, err := h.services.Authorization.ParseToken(headerParts[1])
		if err != nil {
			return errResponse(c, http.StatusUnauthorized, err.Error())
		}

		pkg.InfoPrint("middleware", "ok", "userid: ", userID)
		// Добавляю ID пользователя в контекст
		c.Set(userCtx, userID)

		return next(c)
	}
}

func getUserID(c echo.Context) (int, error) {
	id := c.Get(userCtx)

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}
	return idInt, nil
}
