package handler

// // @Summary SignUp
// // @Tags Auth
// // @Description create account
// // @ID create-account
// // @Accept  json
// // @Produce  json
// // @Param input body core.User true "account info"
// // @Success 200 {integer} integer 1
// // @Failure 400,404 {object} errorResponse
// // @Failure 500 {object} errorResponse
// // @Failure default {object} errorResponse
// // @Router /auth/sign-up [post]
// func (h *Handler) SignUp() echo.HandlerFunc {

// 	return func(c echo.Context) error {
// 		var input core.User

// 		if err := c.Bind(&input); err != nil {
// 			return errResponse(c, http.StatusBadRequest, "invalid input body.")
// 		}

// 		id, err := h.services.Authorization.CreateUser(input)
// 		if err != nil {
// 			return errResponse(c, http.StatusInternalServerError, err.Error())
// 		}

// 		return c.JSON(http.StatusOK, map[string]interface{}{
// 			"id": id,
// 		})

// 	}
// }

// type signInInput struct {
// 	Username string `json:"username" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// }

// // @Summary SignIn
// // @Tags Auth
// // @Description login
// // @ID login
// // @Accept  json
// // @Produce  json
// // @Param input body signInInput true "credentials"
// // @Success 200 {string} string "token"
// // @Failure 400,404 {object} errorResponse
// // @Failure 500 {object} errorResponse
// // @Failure default {object} errorResponse
// // @Router /auth/sign-in [post]
// func (h *Handler) SignIn() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var input signInInput
// 		if err := c.Bind(&input); err != nil {
// 			return errResponse(c, http.StatusBadRequest, "invalid input body")
// 		}

// 		token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
// 		if err != nil {
// 			return errResponse(c, http.StatusInternalServerError, err.Error())
// 		}
// 		return c.JSON(http.StatusOK, map[string]interface{}{
// 			"token": token,
// 		})
// 	}
// }
