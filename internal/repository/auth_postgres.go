package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/sirupsen/logrus"
)

//AuthPostgres реализует логику авторизации и аутентификации
type AuthPostgres struct {
	db *sqlx.DB
}

//NewAuthPostgres конструктор
func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

//Тестирование доступности слоя repository
func (r *AuthPostgres) TestDB() {
	logrus.Info("Info from DB layer")
}

//CreateUser создание пользователя
func (r *AuthPostgres) CreateUser(user core.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (login, password_hash) values ($1, $2) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Login, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

//GetUser получить пользователя из базы
func (r *AuthPostgres) GetUser(login, password string) (core.User, error) {
	// найденный пользователь, парсится в обьект структуры, далее он возвращается на уровень выше
	var user core.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, login, password)
	return user, err
}
