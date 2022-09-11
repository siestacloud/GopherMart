package service

import (
	"github.com/siestacloud/gopherMart/internal/repository"
)

// OrderService реализация бизнес логики обработки номера заказа
type BalanceService struct {
	repo repository.Balance
}

//NewOrderService конструктор
func NewBalanceService(repo repository.Balance) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}

// Create создание нового баланса
// *(используется при авторизации нового пользователя)
func (o *BalanceService) Create(userID int) error {
	return o.repo.Create(userID)
}
