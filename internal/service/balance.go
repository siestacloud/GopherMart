package service

import (
	"github.com/siestacloud/gopherMart/internal/core"
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
func (b *BalanceService) Create(userID int) error {
	return b.repo.Create(userID)
}

// Update обновление текущего количества баллов клиента
func (b *BalanceService) UpdateCurrent(userID int, order *core.Order) error {
	userBalance, err := b.repo.Get(userID)
	if err != nil {
		return err
	}
	userBalance.Current += order.Accrual
	return b.repo.UpdateCurrent(userBalance)
}
