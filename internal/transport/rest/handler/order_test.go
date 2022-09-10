package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/internal/service"
	service_mocks "github.com/siestacloud/gopherMart/internal/service/mocks"
	"github.com/siestacloud/gopherMart/pkg"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateOrder(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockOrder, UserID int, order core.Order)
	// тестовая таблица
	tests := []struct {
		name                 string       // * имя теста
		userID               int          // * уникальный id клиента (вытаскивается из jwt-токена в middleware)
		inputBody            string       // * тело запроса
		inputOrder           core.Order   // * структура пользователя (который передается в метод сервиса)
		mockBehavior         mockBehavior // * функция
		expectedStatusCode   int          // * ожидаемый статус код
		expectedResponseBody string       // * ожидаемое тело ответа
	}{
		//  тест кейсы
		// ! проверка позитивного сценария
		{

			// * новый номер заказа принят в обработку
			name:      "Order accepted",
			userID:    1,
			inputBody: `4561261212345467`,
			inputOrder: core.Order{
				ID: "4561261212345467",
			},
			mockBehavior: func(r *service_mocks.MockOrder, UserID int, order core.Order) {
				r.EXPECT().Create(UserID, order).Return(nil)
			},
			expectedStatusCode:   202,
			expectedResponseBody: ``,
		},
		{
			// * номер заказа уже был загружен этим клиентом
			name:       "Ok",
			userID:     1,
			inputBody:  `10000000`,
			inputOrder: core.Order{ID: "10000000"},
			mockBehavior: func(r *service_mocks.MockOrder, UserID int, order core.Order) {
				r.EXPECT().Create(UserID, order).Return(errors.New("user already have order"))
			},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		// ! проверка негативных сценариев
		{
			// * неверный формат номера заказа
			name:                 "empty client body",
			userID:               1,
			inputBody:            ``,
			inputOrder:           core.Order{ID: ""},
			mockBehavior:         func(r *service_mocks.MockOrder, UserID int, order core.Order) {},
			expectedStatusCode:   422,
			expectedResponseBody: `{"message":"order format failure"}` + "\n",
		},
		{
			// * неверный формат номера заказа
			name:                 "some text in client body",
			userID:               1,
			inputBody:            `hello I,m hacking your order system!!!`,
			inputOrder:           core.Order{ID: ""},
			mockBehavior:         func(r *service_mocks.MockOrder, UserID int, order core.Order) {},
			expectedStatusCode:   422,
			expectedResponseBody: `{"message":"order format failure"}` + "\n",
		},
		{
			// * неверный формат номера заказа
			name:                 `json in client body`,
			userID:               717,
			inputBody:            `{"hackJSON": "//0213ddsd2/dsd","0_0": "HackScripts[<djsldnas><>]"}`,
			inputOrder:           core.Order{ID: ""},
			mockBehavior:         func(r *service_mocks.MockOrder, UserID int, order core.Order) {},
			expectedStatusCode:   422,
			expectedResponseBody: `{"message":"order format failure"}` + "\n",
		},
		{
			// * неверный формат запроса
			name:                 "0s client body",
			userID:               717,
			inputBody:            `000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000`,
			inputOrder:           core.Order{ID: ""},
			mockBehavior:         func(r *service_mocks.MockOrder, UserID int, order core.Order) {},
			expectedStatusCode:   422,
			expectedResponseBody: `{"message":"order format failure"}` + "\n",
		},
		{
			// * номер заказа уже был загружен другим клиентом
			name:       "Order another user",
			userID:     1,
			inputBody:  `4561261212345467`,
			inputOrder: core.Order{ID: "4561261212345467"},
			mockBehavior: func(r *service_mocks.MockOrder, UserID int, order core.Order) {
				r.EXPECT().Create(UserID, order).Return(errors.New("another user order"))
			},
			expectedStatusCode:   409,
			expectedResponseBody: `{"message":"another user order"}` + "\n",
		},
		{
			// * внутренняя ошибка сервера
			name:       "Internal server error",
			userID:     1,
			inputBody:  `4561261212345467`,
			inputOrder: core.Order{ID: "4561261212345467"},
			mockBehavior: func(r *service_mocks.MockOrder, UserID int, order core.Order) {
				r.EXPECT().Create(UserID, order).Return(errors.New("some err in service or repository layers..."))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"internal server error"}` + "\n",
		},
	}

	// В цикле итерируемся по тестовой таблице
	for _, test := range tests {

		// * вызываем метод RUN у объекта t )
		// * передаем имя теста и функцию
		// * тесты запускаются параллельно в отдельных горутинах
		t.Run(test.name, func(t *testing.T) {

			// в теле тест функции инициализируем зависимости
			// * создаем контроллер мока слоя сервис
			// * вызываем метод finish (оссобенность библиотеки
			// * для каждого теста нужно создавать контроллер и финишировать его по выполнению теста)
			c := gomock.NewController(t)
			defer c.Finish()

			// * создаем мок слоя сервис, передаем контроллер как аргумент
			order := service_mocks.NewMockOrder(c)

			// * в данном тестовом сценарии ожидаем получить
			// * вызов метода сервиса и получить в качестве аргумента данную структуру
			test.mockBehavior(order, test.userID, test.inputOrder)

			// * инициализируем слой service, имплементируем интерфейс Authorization моком auth
			services := &service.Service{Order: order}
			handler := NewHandler(services)

			// * инициализация тестового ендпоитна
			e := echo.New()

			e.Validator = pkg.NewCustomValidator(validator.New())

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			s := e.NewContext(req, rec)
			s.Set(userCtx, test.userID)

			q := handler.CreateOrder()
			// * Проверка корректности JWT-токена в заголовке и статуса ответа
			if assert.NoError(t, q(s)) {
				assert.Equal(t, test.expectedStatusCode, rec.Code)
				assert.Equal(t, test.expectedResponseBody, rec.Body.String())

			}
		})
	}
}
