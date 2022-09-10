package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/siestacloud/gopherMart/internal/config"
	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/internal/repository"
	repository_mocks "github.com/siestacloud/gopherMart/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateUser(t *testing.T) {
	type mockBehavior func(r *repository_mocks.MockAuthorization, user core.User)
	// тестовая таблица
	tests := []struct {
		name           string       // * имя теста
		inputUser      core.User    // * структура пользователя (который передается в метод сервиса)
		hashUser       core.User    // * структура пользователя (который передается в метод сервиса)
		mockBehavior   mockBehavior // * функция
		expectedUserID int
	}{
		// тест кейсы
		{
			name: "OK",
			inputUser: core.User{
				Login:    "poul2",
				Password: "pass",
			},
			hashUser: core.User{
				Login:    "poul2",
				Password: "686a7172686a7177313234363137616a6668616a739d4e1e23bd5b727046a9e3b4b7db57bd8d6ee684",
			},
			// * создаем поведение для объекта мока
			// * указываем что ожидаем получить вызов CreateUser именно с той структурой user которую передаем
			// * и этот метод должен вернуть ID и nil в качестве ошибки
			mockBehavior: func(r *repository_mocks.MockAuthorization, user core.User) {
				r.EXPECT().CreateUser(user).Return(1, nil)
			},
			// * указываем ожидаемый статус код и тело ответа
			expectedUserID: 1,
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
			auth := repository_mocks.NewMockAuthorization(c)

			// * в данном тестовом сценарии ожидаем получить
			// * вызов метода сервиса и получить в качестве аргумента данную структуру пользователя
			test.mockBehavior(auth, test.hashUser)

			// * инициализируем слой service, имплементируем интерфейс Authorization моком auth
			repository := &repository.Repository{Authorization: auth}

			service := NewService(&config.Cfg{}, repository)

			res, err := service.CreateUser(test.inputUser)
			if err != nil {
				assert.Error(t, err)
			}
			// * Проверка корректности JWT-токена в заголовке и статуса ответа
			assert.Equal(t, test.expectedUserID, res)
		})
	}
}

func TestService_GenerateToken(t *testing.T) {
	type mockBehavior func(r *repository_mocks.MockAuthorization, login, password string)
	// тестовая таблица
	tests := []struct {
		name          string       // * имя теста
		login         string       // * структура пользователя (который передается в метод сервиса)
		password      string       // * структура пользователя (который передается в метод сервиса)
		mockBehavior  mockBehavior // * функция
		expectedToken string
	}{
		// тест кейсы
		{
			name:     "OK",
			login:    "user",
			password: "qwerty@@!#333",
			// * создаем поведение для объекта мока
			// * указываем что ожидаем получить вызов CreateUser именно с той структурой user которую передаем
			// * и этот метод должен вернуть ID и nil в качестве ошибки
			mockBehavior: func(r *repository_mocks.MockAuthorization, login, password string) {
				r.EXPECT().GetUser(login, generatePasswordHash(password)).Return(&core.User{ID: 1}, nil)
			},
			// * указываем ожидаемый статус код и тело ответа
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			name:     "OK",
			login:    "somelogin",
			password: "passWD@@#3321",
			// * создаем поведение для объекта мока
			// * указываем что ожидаем получить вызов CreateUser именно с той структурой user которую передаем
			// * и этот метод должен вернуть ID и nil в качестве ошибки
			mockBehavior: func(r *repository_mocks.MockAuthorization, login, password string) {
				r.EXPECT().GetUser(login, generatePasswordHash(password)).Return(&core.User{ID: 998}, nil)
			},
			// * указываем ожидаемый статус код и тело ответа
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
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
			auth := repository_mocks.NewMockAuthorization(c)

			// * в данном тестовом сценарии ожидаем получить
			// * вызов метода сервиса и получить в качестве аргумента данную структуру пользователя
			test.mockBehavior(auth, test.login, test.password)

			// * инициализируем слой service, имплементируем интерфейс Authorization моком auth
			repository := &repository.Repository{Authorization: auth}

			service := NewService(&config.Cfg{}, repository)

			res, err := service.GenerateToken(test.login, test.password)
			if err != nil {
				assert.Error(t, err)
			}
			fmt.Println("RES ", res)
			// * Проверка корректности JWT-токена
			assert.Equal(t, test.expectedToken, strings.Split(res, ".")[0])
		})
	}
}

// func TestService_ParseToken(t *testing.T) {
// 	// тестовая таблица
// 	tests := []struct {
// 		name           string // * имя теста
// 		accessToken    string
// 		expectedUserID int
// 	}{
// 		// тест кейсы
// 		{
// 			name:        "OK",
// 			accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjI1MzE4NDgsImlhdCI6MTY2MjQ4ODY0OCwidXNlcl9pZCI6MX0.nBAA1CAj_ijpZ06VksxaI-um7pH7RXKxoW0xbjfGVpc",
// 			// * указываем ожидаемый статус код и тело ответа
// 			expectedUserID: 1,
// 		},
// 		{
// 			name:        "OK",
// 			accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjI1MzIxNzIsImlhdCI6MTY2MjQ4ODk3MiwidXNlcl9pZCI6MjB9.YFqiuWku_OoUzCsJ35OZvcy2vlS97hILDWROrsEnC_w",
// 			// * указываем ожидаемый статус код и тело ответа
// 			expectedUserID: 20,
// 		},
// 		{
// 			name:        "OK",
// 			accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjI1MzIyMTgsImlhdCI6MTY2MjQ4OTAxOCwidXNlcl9pZCI6OTk4fQ.a70jOsxNRxhCKOsTbAqdQUjcw_3DwUFrTIWLeVUPXRg",
// 			// * указываем ожидаемый статус код и тело ответа
// 			expectedUserID: 998,
// 		},
// 	}

// 	// В цикле итерируемся по тестовой таблице
// 	for _, test := range tests {

// 		// * вызываем метод RUN у объекта t )
// 		// * передаем имя теста и функцию
// 		// * тесты запускаются параллельно в отдельных горутинах
// 		t.Run(test.name, func(t *testing.T) {

// 			// в теле тест функции инициализируем зависимости
// 			// * создаем контроллер мока слоя сервис
// 			// * вызываем метод finish (оссобенность библиотеки
// 			// * для каждого теста нужно создавать контроллер и финишировать его по выполнению теста)
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			// * создаем мок слоя сервис, передаем контроллер как аргумент
// 			auth := repository_mocks.NewMockAuthorization(c)

// 			// * инициализируем слой service, имплементируем интерфейс Authorization моком auth
// 			repository := &repository.Repository{Authorization: auth}

// 			service := NewService(&config.Cfg{}, repository)

// 			res, err := service.ParseToken(test.accessToken)
// 			if err != nil {
// 				assert.Error(t, err)
// 			}
// 			// * Проверка корректности JWT-токена
// 			assert.Equal(t, test.expectedUserID, res)
// 		})
// 	}
// }

func TestAuthService_ParseToken(t *testing.T) {
	type fields struct {
		repo repository.Authorization
	}
	type args struct {
		accessToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthService{
				repo: tt.fields.repo,
			}
			got, err := s.ParseToken(tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthService.ParseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
