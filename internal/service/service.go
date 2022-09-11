package service

import (
	"github.com/siestacloud/gopherMart/internal/core"

	"github.com/siestacloud/gopherMart/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	Test()
	CreateUser(user core.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Order interface {
	Create(userID int, order core.Order) error
	GetUserByOrder(orderID string) (int, error)
	GetListOrders(userID int) ([]core.Order, error)
}

type Accrual interface {
	GetOrderAccrual(order *core.Order) error
}

type Balance interface {
	Create(userId int) error
	Get(userID int) (*core.Balance, error)
	UpdateCurrent(userID int, order *core.Order) error
}

// Главный тип слоя SVC, который встраивается как зависимость в слое TRANSPORT
type Service struct {
	Authorization
	Accrual
	Balance
	Order
}

// Конструктор слоя SVC
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Accrual:       NewAccrualService(repos.Accrual),
		Balance:       NewBalanceService(repos.Balance),
		Order:         NewOrderService(repos.Order),
	}
}
