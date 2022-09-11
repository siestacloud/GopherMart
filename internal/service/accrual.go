package service

import (
	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/internal/repository"
)

// OrderService реализация бизнес логики обработки номера заказа
type AccrualService struct {
	repo repository.Accrual
}

//NewOrderService конструктор
func NewAccrualService(repo repository.Accrual) *AccrualService {
	return &AccrualService{
		repo: repo,
	}
}

func (o *AccrualService) GetOrderAccrual(order *core.Order) error {
	return o.repo.GetOrderAccrual(order)
}
