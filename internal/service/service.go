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
	Create(userId int, order core.Order) error
	GetUserByOrder(order int) (int, error)
	GetListOrders(userID int) ([]core.Order, error)
}

// Главный тип слоя SVC, который встраивается как зависимость в слое TRANSPORT
type Service struct {
	Authorization
	Order
}

// Конструктор слоя SVC
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Order:         NewOrderService(repos.Order),
	}
}
