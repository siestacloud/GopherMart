package service

import (
	"errors"
	"fmt"

	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/internal/repository"
	"github.com/siestacloud/gopherMart/pkg"
)

// OrderService реализация бизнес логики обработки номера заказа
type OrderService struct {
	repo repository.Order
}

//NewOrderService конструктор
func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

//Create проверка номера заказа по алгоритму ЛУНА и сохранение его в базе (с привязкой к конкретному пользователю)
func (o *OrderService) Create(userID int, order core.Order) error {

	if !pkg.Valid(order.ID) {
		pkg.WarnPrint("service", "lune alg err", "invalid order")
		return errors.New("lune alg invalid order")
	}

	userDB, err := o.repo.GetUserByOrder(order.ID)
	if err != nil {
		pkg.WarnPrint("service", "get user by order", err)
		if err = o.repo.Create(userID, order); err != nil {
			return err
		}
		return err
	}
	fmt.Println(userDB, "==", userID)
	if userDB == userID {
		return errors.New("user already have order")
	}

	return errors.New("another user order")

}

func (o *OrderService) GetUserByOrder(order int) (int, error) {
	return o.repo.GetUserByOrder(order)
}
