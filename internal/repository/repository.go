package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/siestacloud/gopherMart/internal/config"
	"github.com/siestacloud/gopherMart/internal/core"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Authorization interface {
	TestDB()
	CreateUser(user core.User) (int, error)
	GetUser(username, password string) (*core.User, error)
}

type Order interface {
	Create(userId int, order core.Order, createTime string) error
	GetUserByOrder(orderID string) (int, error)
	GetListOrders(userID int) ([]core.Order, error)
}

type Accrual interface {
	GetOrderAccrual(order *core.Order) error
}

// Главный тип слоя repository, который встраивается как зависимость в слое SVC
type Repository struct {
	Authorization
	Accrual
	Order
}

//Конструктор слоя repository
func NewRepository(cfg *config.Cfg, db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Accrual:       NewAccrualAPI(cfg),
		Order:         NewOrderPostgres(db),
	}
}
