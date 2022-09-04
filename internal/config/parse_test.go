package config_test

import (
	"testing"

	"github.com/siestacloud/gopherMart/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct { // добавился слайс тестов
		name   string
		values config.Cfg
		want   config.Cfg
	}{
		{
			name:   "test #1",    // описывается каждый тест
			values: config.Cfg{}, // значения, которые будет принимать функция
			want: config.Cfg{
				Server: config.Server{
					Address:              "localhost:8080",
					URLPostgres:          "not set",
					AccrualSystemAddress: "not set",
					Logrus: config.Logrus{
						LogLevel: "info",
						JSON:     false},
				},
			}, // ожидаемое значение
		},
	}
	for _, tt := range tests { // цикл по всем тестам
		t.Run(tt.name, func(t *testing.T) {

			if err := config.Parse(&tt.values); err != nil {
				t.Errorf(" want %v", tt.want)
			}
			assert.Equal(t, tt.want, tt.values)
		})
	}
}

// func TestFullName(t *testing.T) {

// 	tests := []struct { // добавился слайс тестов
// 		name   string
// 		values User
// 		want   string
// 	}{
// 		{
// 			name:   "simple test #1",                            // описывается каждый тест
// 			values: User{FirstName: "Poul", LastName: "Siesta"}, // значения, которые будет принимать функция
// 			want:   "Poul Siesta",                               // ожидаемое значение
// 		},
// 		{
// 			name:   "test with empty",
// 			values: User{},
// 			want:   " ",
// 		},
// 	}
// 	for _, tt := range tests { // цикл по всем тестам
// 		t.Run(tt.name, func(t *testing.T) {
// 			if sum := tt.values.FullName(); sum != tt.want {
// 				t.Errorf("FullName = %v, want %v", sum, tt.want)
// 			}

// 			//
// 			//or
// 			//

// 			assert.Equal(t, tt.want, tt.values.FullName())
// 		})
// 	}
// }
