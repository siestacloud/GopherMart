package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/siestacloud/gopherMart/internal/core"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Authorization interface {
	TestDB()
	CreateUser(user core.User) (int, error)
	GetUser(username, password string) (*core.User, error)
}

type Order interface {
	Create(userId int, order core.Order, status, createTime string) error
	GetUserByOrder(orderID string) (int, error)
	GetListOrders(userID int) ([]core.Order, error)
}

// Главный тип слоя repository, который встраивается как зависимость в слое SVC
type Repository struct {
	Authorization
	Order
}

//Конструктор слоя repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Order:         NewOrderPostgres(db),
	}
}
