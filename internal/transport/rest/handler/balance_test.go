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

// TestHandler_GetBalance тестирование логики обработки GET запроса от клиента для: получения текущего баланса клиента;
func TestHandler_GetBalance(t *testing.T) {
	type mockBehavior func(b *service_mocks.MockBalance, UserID int)
	// тестовая таблица
	tests := []struct {
		name                 string       // * имя теста
		userID               int          // * уникальный id клиента (вытаскивается из jwt-токена в middleware)
		mockBehavior         mockBehavior // * функция
		expectedStatusCode   int          // * ожидаемый статус код
		expectedResponseBody string       // * ожидаемое тело ответа
	}{
		//  тест кейсы
		// ! проверка позитивного сценария
		{
			// * новый номер заказа принят в обработку
			name:   "get balance ok",
			userID: 1,
			mockBehavior: func(b *service_mocks.MockBalance, UserID int) {
				b.EXPECT().Get(UserID).Return(&core.Balance{Current: 100, Withdrawn: 400}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"ID":0,"current":100,"withdrawn":400}` + "\n",
		},
		// ! проверка негативных сценариев
		{
			// * обработка внутренней ошибки
			name:   "some internal error",
			userID: 1,
			mockBehavior: func(b *service_mocks.MockBalance, UserID int) {
				b.EXPECT().Get(UserID).Return(nil, errors.New("some internal error..."))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"internal server error"}` + "\n",
		},
	}

	// В цикле итерируемся по тестовой таблице
	for _, test := range tests {

		// * вызываем метод RUN у объекта t)
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
			balance := service_mocks.NewMockBalance(c)

			// * в данном тестовом сценарии ожидаем получить
			// * вызов метода сервиса и получить
			test.mockBehavior(balance, test.userID)

			// * инициализируем слой service, имплементируем интерфейс Balance
			services := &service.Service{Balance: balance}
			handler := NewHandler(services)

			// * инициализация тестового ендпоитна
			e := echo.New()

			e.Validator = pkg.NewCustomValidator(validator.New())

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
			s := e.NewContext(req, rec)
			s.Set(userCtx, test.userID)

			q := handler.GetBalance()
			if assert.NoError(t, q(s)) {
				assert.Equal(t, test.expectedStatusCode, rec.Code)
				assert.Equal(t, test.expectedResponseBody, rec.Body.String())
			}
		})
	}
}

// TestHandler_WithdrawBalance тестирование логики обработки POST запроса от клиента для: списания баллов;
func TestHandler_WithdrawBalance(t *testing.T) {
	type mockBehavior func(b *service_mocks.MockBalance, o *service_mocks.MockOrder, UserID int, order core.Order)
	// тестовая таблица
	tests := []struct {
		name                 string // * имя теста
		userID               int    // * уникальный id клиента (вытаскивается из jwt-токена в middleware)
		inputBody            string
		inputOrder           core.Order
		mockBehavior         mockBehavior // * функция
		expectedStatusCode   int          // * ожидаемый статус код
		expectedResponseBody string       // * ожидаемое тело ответа
	}{
		//  тест кейсы
		// ! проверка позитивного сценария
		{
			// * новый номер заказа принят в обработку
			name:       "get balance ok",
			userID:     1,
			inputBody:  `{"order":"2377225624", "sum":751}`,
			inputOrder: core.Order{Number: "2377225624", Sum: 751},
			mockBehavior: func(b *service_mocks.MockBalance, o *service_mocks.MockOrder, UserID int, order core.Order) {
				o.EXPECT().Create(UserID, order).Return(nil)
				b.EXPECT().Withdrawal(UserID, order.Sum).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		// ! проверка негативных сценариев
		// todo добавить тесты
	}

	// В цикле итерируемся по тестовой таблице
	for _, test := range tests {

		// * вызываем метод RUN у объекта t)
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
			balance := service_mocks.NewMockBalance(c)
			order := service_mocks.NewMockOrder(c)

			// * в данном тестовом сценарии ожидаем получить
			// * вызов метода сервиса и получить
			test.mockBehavior(balance, order, test.userID, test.inputOrder)

			// * инициализируем слой service, имплементируем интерфейс Balance
			services := &service.Service{Balance: balance, Order: order}
			handler := NewHandler(services)

			// * инициализация тестового ендпоитна
			e := echo.New()

			e.Validator = pkg.NewCustomValidator(validator.New())

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			s := e.NewContext(req, rec)
			s.Set(userCtx, test.userID)

			q := handler.WithdrawBalance()
			if assert.NoError(t, q(s)) {
				assert.Equal(t, test.expectedStatusCode, rec.Code)
				assert.Equal(t, test.expectedResponseBody, rec.Body.String())
			}
		})
	}
}

// TestHandler_WithdrawalsBalance тестирование логики обработки GET запроса от клиента для: Получение информации о выводе средств;
func TestHandler_WithdrawalsBalance(t *testing.T) {
	type mockBehavior func(o *service_mocks.MockOrder, UserID int)
	// тестовая таблица
	tests := []struct {
		name                 string       // * имя теста
		userID               int          // * уникальный id клиента (вытаскивается из jwt-токена в middleware)
		mockBehavior         mockBehavior // * функция
		expectedStatusCode   int          // * ожидаемый статус код
		expectedResponseBody string       // * ожидаемое тело ответа
	}{
		//  тест кейсы
		// ! проверка позитивного сценария
		{
			// * получение информации о выводе средств
			name:   "get balance ok",
			userID: 1,
			mockBehavior: func(o *service_mocks.MockOrder, UserID int) {
				o.EXPECT().GetListOrders(UserID).Return([]core.Order{{Number: "2377225624", Sum: 210, CreateTime: "2022-09-14T15:21:16+03:00"}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"order":"2377225624","sum":210,"processed_at":"2022-09-14T15:21:16+03:00"}]` + "\n",
		},
		// ! проверка негативных сценариев

	}

	// В цикле итерируемся по тестовой таблице
	for _, test := range tests {

		// * вызываем метод RUN у объекта t)
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
			// * вызов метода сервиса и получить
			test.mockBehavior(order, test.userID)

			// * инициализируем слой service, имплементируем интерфейс Balance
			services := &service.Service{Order: order}
			handler := NewHandler(services)

			// * инициализация тестового ендпоитна
			e := echo.New()

			e.Validator = pkg.NewCustomValidator(validator.New())

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
			s := e.NewContext(req, rec)
			s.Set(userCtx, test.userID)

			q := handler.WithdrawalsBalance()
			if assert.NoError(t, q(s)) {
				assert.Equal(t, test.expectedStatusCode, rec.Code)
				assert.Equal(t, test.expectedResponseBody, rec.Body.String())
			}
		})
	}
}
