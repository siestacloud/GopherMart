package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/siestacloud/gopherMart/internal/service"
	service_mocks "github.com/siestacloud/gopherMart/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

// TestHandler_userIdentity тесты для Middleware
// * тестирование логики добавления в контекст запроса ID клиента из заголовка Authorization)
func TestHandler_userIdentity(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *service_mocks.MockAuthorization, token string)

	testTable := []struct {
		name                 string       // * имя теста
		headerName           string       // * имя заголовка
		headerValue          string       // * значение заголовка
		token                string       // * токен
		mockBehavior         mockBehavior // * функция заглушка
		expectedStatusCode   int          // * ожидаемый статус код ответа
		expectedResponseBody string       // * ожидаемое тело ответа (userid from token)
	}{
		// ! проверка позитивных сценариев
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *service_mocks.MockAuthorization, token string) {
				r.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1" + "\n",
		},
		// ! проверка негативных сценариев
		{ // * отсутствует header Authorization
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *service_mocks.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}` + "\n",
		},
		{ // * неправильное значение header
			name:                 "Invalid Header Value",
			headerName:           "Authorization",
			headerValue:          "Bearr token",
			token:                "token",
			mockBehavior:         func(r *service_mocks.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}` + "\n",
		},
		{ // * отсутствует токен
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "",
			mockBehavior:         func(r *service_mocks.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"token is empty"}` + "\n",
		},
		{ // * неправильный токен (в тестовом сценарии мок функция возвращает ошибку)
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *service_mocks.MockAuthorization, token string) {
				r.EXPECT().ParseToken(token).Return(0, errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid token"}` + "\n",
		},
	}

	// В цикле итерируемся по тестовой таблице
	for _, test := range testTable {

		// * вызываем метод RUN у объекта t )
		// * передаем имя теста и функцию
		// * тесты запускаются параллельно в отдельных горутинах
		t.Run(test.name, func(t *testing.T) {

			// в теле тест функции инициализирую зависимости
			// * создаем контроллер мока слоя сервис
			// * вызываем метод finish (оссобенность библиотеки
			// * для каждого теста нужно создавать контроллер и финишировать его по выполнению теста)
			c := gomock.NewController(t)
			defer c.Finish()

			// * создаем мок слоя сервис, передаем контроллер как аргумент
			auth := service_mocks.NewMockAuthorization(c)

			// * в данном тестовом сценарии ожидаем получить
			// * вызов метода сервиса и получить в качестве аргумента токен из заголовка
			test.mockBehavior(auth, test.token)

			// * инициализируем слой service, имплементирую интерфейс Authorization моком auth
			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// * инициализация тестового ендпоитна
			e := echo.New()
			e.POST("/", func(c echo.Context) error {
				id := c.Get(userCtx).(int)
				return c.JSON(http.StatusOK, id)
			})
			e.Use(handler.UserIdentity) //  тестируемая middleware (JWT token auth)

			// * формирую новый post запрос от клиента
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(test.headerName, test.headerValue)
			e.ServeHTTP(rec, req)

			// * Проверка корректности JWT-токена в заголовке и статуса ответа
			assert.Equal(t, test.expectedStatusCode, rec.Code)
			assert.Equal(t, test.expectedResponseBody, rec.Body.String())

		})
	}
}

// TestGetUserID
// * тестирование логики получения уникального ID клиента из контекста запроса
func TestGetUserID(t *testing.T) {
	// * getContext возвращает контекст для (используется в тестовых сценариях)
	var getContext = func(id interface{}) *echo.Context {
		e := echo.New()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := e.NewContext(req, rec)
		ctx.Set(userCtx, id)
		return &ctx
	}

	// * тестовые сценарии
	testTable := []struct {
		name       string        // * имя теста
		ctx        *echo.Context // * контекст
		shouldFail bool
		expectedID int // * ID клиента
	}{
		{
			name:       "Ok",
			ctx:        getContext(1),
			expectedID: 1,
			shouldFail: false,
		},
		{
			name:       "Empty",
			ctx:        getContext("invalid_user_id_pcjwpwx"),
			shouldFail: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			id, err := getUserID(*test.ctx)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, id, test.expectedID)
		})
	}
}
