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

func TestHandler_SignUp(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockAuthorization, user core.User)
	// тестовая таблица
	tests := []struct {
		name                 string       // * имя теста
		inputBody            string       // * тело запроса
		inputUser            core.User    // * структура пользователя (который передается в метод сервиса)
		mockBehavior         mockBehavior // * функция
		expectedToken        string
		expectedStatusCode   int    // * ожидаемый статус код
		expectedResponseBody string // * ожидаемое тело ответа
	}{
		// тест кейсы
		// ! проверка позитивного сценария
		{
			name:      "OK",
			inputBody: `{"login": "poul2","password": "pass"}`,
			inputUser: core.User{
				Login:    "poul2",
				Password: "pass",
			},
			// * создаем поведение для объекта мока
			// * указываем что ожидаем получить вызов CreateUser именно с той структурой user которую передаем
			// * и этот метод должен вернуть ID и nil в качестве ошибки
			mockBehavior: func(r *service_mocks.MockAuthorization, user core.User) {
				r.EXPECT().CreateUser(user).Return(1, nil)
				r.EXPECT().GenerateToken(user.Login, user.Password).Return("qazwsxedC", nil)
			},
			expectedToken: "Bearer qazwsxedC",
			// * указываем ожидаемый статус код и тело ответа
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		// ! проверка негативных сценариев
		{
			// * логин занят
			name:      "login busy",
			inputBody: `{"login": "poul2","password": "pass"}`,
			inputUser: core.User{
				Login:    "poul2",
				Password: "pass",
			},
			mockBehavior: func(r *service_mocks.MockAuthorization, user core.User) {
				r.EXPECT().CreateUser(user).Return(0, errors.New("login busy"))
			},
			expectedStatusCode:   409,
			expectedResponseBody: `{"message":"login busy"}` + "\n",
		},

		{
			// * поля в client body пустые
			name:                 "Empty fields",
			inputBody:            `{"login": "","password": ""}`,
			inputUser:            core.User{},
			mockBehavior:         func(r *service_mocks.MockAuthorization, user core.User) {}, // * вызывается в этом тесте
			expectedToken:        "",
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"validate failure"}` + "\n",
		},
		{
			// * поле в client body login пустое
			name:                 "Empty client field",
			inputBody:            `{"login": "","password": "passw@@ord"}`,
			inputUser:            core.User{},
			mockBehavior:         func(r *service_mocks.MockAuthorization, user core.User) {}, // * не вызывается в этом тесте
			expectedToken:        "",
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"validate failure"}` + "\n",
		},
		{
			// * поле в client body password пустое
			name:                 "Empty password field",
			inputBody:            `{"login": "","password": ""}`,
			inputUser:            core.User{},
			mockBehavior:         func(r *service_mocks.MockAuthorization, user core.User) {}, // * не вызывается в этом тесте
			expectedToken:        "",
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"validate failure"}` + "\n",
		},
		{
			// * поля в client body отсутствуют
			name:                 "No fields in client body request",
			inputBody:            `{"another1": "someval1","another2": "someval2"}`,
			inputUser:            core.User{},
			mockBehavior:         func(r *service_mocks.MockAuthorization, user core.User) {}, // * не вызывается в этом тесте
			expectedToken:        "",
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"validate failure"}` + "\n",
		},
		{ // * внутренняя ошибка сервера
			name:      "internal server error",
			inputBody: `{"login": "poul2","password": "pass"}`,
			inputUser: core.User{
				Login:    "poul2",
				Password: "pass",
			},
			// * создаем поведение для объекта мока
			// * указываем что ожидаем получить вызов CreateUser именно с той структурой user которую передаем
			// * и этот метод должен вернуть ID и nil в качестве ошибки
			mockBehavior: func(r *service_mocks.MockAuthorization, user core.User) {
				r.EXPECT().CreateUser(user).Return(0, errors.New("some internal error"))
			},
			// * указываем ожидаемый статус код и тело ответа
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
			auth := service_mocks.NewMockAuthorization(c)

			// * в данном тестовом сценарии ожидаем получить
			// * вызов метода сервиса и получить в качестве аргумента данную структуру пользователя
			test.mockBehavior(auth, test.inputUser)

			// * инициализируем слой service, имплементируем интерфейс Authorization моком auth
			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// * инициализация тестового ендпоитна
			e := echo.New()
			e.Validator = pkg.NewCustomValidator(validator.New())

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			s := e.NewContext(req, rec)

			q := handler.SignUp()
			// * Проверка корректности JWT-токена в заголовке и статуса ответа
			if assert.NoError(t, q(s)) {
				assert.Equal(t, test.expectedStatusCode, rec.Code)
				assert.Equal(t, test.expectedToken, rec.Header().Get("authorization"))
				assert.Equal(t, test.expectedResponseBody, rec.Body.String())

			}
		})
	}
}

func TestHandler_SignIn(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockAuthorization, user core.User)
	// тестовая таблица
	tests := []struct {
		name                 string       // * имя теста
		inputBody            string       // * тело запроса
		inputUser            core.User    // * структура пользователя (который передается в метод сервиса)
		mockBehavior         mockBehavior // * функция
		expectedToken        string
		expectedStatusCode   int    // * ожидаемый статус код
		expectedResponseBody string // * ожидаемое тело ответа
	}{
		// ! проверка позитивных сценариев
		// - пользователь успешно аутентифицирован;
		{
			name:      "OK",
			inputBody: `{"login": "poul2","password": "pass"}`,
			inputUser: core.User{
				Login:    "poul2",
				Password: "pass",
			},
			// * создаем поведение для объекта мока
			// * указываем что ожидаем получить вызов CreateUser именно с той структурой user которую передаем
			// * и этот метод должен вернуть ID и nil в качестве ошибки
			mockBehavior: func(r *service_mocks.MockAuthorization, user core.User) {
				r.EXPECT().GenerateToken(user.Login, user.Password).Return("qazwsxedC", nil)
			},
			expectedToken: "Bearer qazwsxedC",
			// * указываем ожидаемый статус код и тело ответа
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},

		// ! проверка негативных сценариев
		{
			// * поле в client body login пустое
			name:                 "Empty client field",
			inputBody:            `{"login": "","password": "passw@@ord"}`,
			inputUser:            core.User{},
			mockBehavior:         func(r *service_mocks.MockAuthorization, user core.User) {}, // * не вызывается в этом тесте
			expectedToken:        "",
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"validate failure"}` + "\n",
		},
		{
			// * поле в client body password пустое
			name:                 "Empty password field",
			inputBody:            `{"login": "","password": ""}`,
			inputUser:            core.User{},
			mockBehavior:         func(r *service_mocks.MockAuthorization, user core.User) {}, // * не вызывается в этом тесте
			expectedToken:        "",
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"validate failure"}` + "\n",
		},
		{
			// * поля в client body отсутствуют
			name:                 "No fields in client body request",
			inputBody:            `{"another1": "someval1","another2": "someval2"}`,
			inputUser:            core.User{},
			mockBehavior:         func(r *service_mocks.MockAuthorization, user core.User) {}, // * не вызывается в этом тесте
			expectedToken:        "",
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"validate failure"}` + "\n",
		},
		{
			// * неверная пара логин/пароль
			name:      "invalid username/password pair",
			inputBody: `{"login": "poul2","password": "pass"}`,
			inputUser: core.User{
				Login:    "poul2",
				Password: "pass",
			},
			mockBehavior: func(r *service_mocks.MockAuthorization, user core.User) {
				r.EXPECT().GenerateToken(user.Login, user.Password).Return("", errors.New("invalid username/password pair"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid username/password pair"}` + "\n",
		},
		{ // * внутренняя ошибка сервера
			name:      "internal server error",
			inputBody: `{"login": "poul2","password": "pass"}`,
			inputUser: core.User{
				Login:    "poul2",
				Password: "pass",
			},
			mockBehavior: func(r *service_mocks.MockAuthorization, user core.User) {
				r.EXPECT().GenerateToken(user.Login, user.Password).Return("", errors.New("some internal error"))
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
			auth := service_mocks.NewMockAuthorization(c)

			// * в данном тестовом сценарии ожидаем получить
			// * вызов метода сервиса и получить в качестве аргумента данную структуру пользователя
			test.mockBehavior(auth, test.inputUser)

			// * инициализируем слой service, имплементируем интерфейс Authorization моком auth
			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// * инициализация тестового ендпоитна
			e := echo.New()
			e.Validator = pkg.NewCustomValidator(validator.New())

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			s := e.NewContext(req, rec)

			q := handler.SignIn()
			// * Проверка корректности JWT-токена в заголовке и статуса ответа
			if assert.NoError(t, q(s)) {
				assert.Equal(t, test.expectedStatusCode, rec.Code)
				assert.Equal(t, test.expectedToken, rec.Header().Get("authorization"))
				assert.Equal(t, test.expectedResponseBody, rec.Body.String())

			}
		})
	}
}
