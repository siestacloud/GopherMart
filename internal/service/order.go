package service

import (
	"errors"
	"sort"
	"time"

	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/internal/repository"
	"github.com/siestacloud/gopherMart/pkg"
)

var (
	statusNew string = "NEW"
	//// statusProcessing string = "PROCESSING"
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

	userDB, err := o.repo.GetUserByOrder(order.Number) // * попытка определить клиента по номеру заказа
	if err != nil {
		pkg.WarnPrint("service", "get user by order", err)

		order.CreateTime = time.Now().Format(time.RFC3339)
		order.WithdrawnTime = time.Now().Format(time.RFC3339)

		if err = o.repo.Create(userID, order); err != nil { // * клиент c таким номером не был найден, заказ сохраняется в бд
			return err
		}
		return err
	}
	if userDB == userID {
		return errors.New("user already have order")
	}

	return errors.New("another user order")

}

func (o *OrderService) GetUserByOrder(order string) (int, error) {
	return o.repo.GetUserByOrder(order)
}

func (o *OrderService) GetListOrders(userID int) ([]core.Order, error) {
	list, err := o.repo.GetListOrders(userID)
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool {
		ti, _ := time.Parse(time.RFC3339, list[i].CreateTime)
		tj, _ := time.Parse(time.RFC3339, list[j].CreateTime)
		return ti.Before(tj)
	})
	pkg.InfoPrint("service", "OK", list)
	return list, err
}

// //Create проверка номера заказа со списанием по алгоритму ЛУНА и сохранение его в базе (с привязкой к конкретному пользователю)
// func (o *OrderService) Withdraw(userID int, order core.Order) error {

// 	// * проверка номера заказа по алгоритму Луна
// 	if err := pkg.Valid(order.Number); err != nil {
// 		pkg.WarnPrint("service", "lune alg err", err)
// 		return err
// 	}

// 	userDB, err := o.repo.GetUserByOrder(order.Number) // * попытка определить клиента по номеру заказа
// 	if err != nil {
// 		pkg.WarnPrint("service", "order doesn't belong any user", err)
// 		order.WithdrawnTime = time.Now().Format(time.RFC3339)

// 		if err = o.repo.Create(userID, order); err != nil { // * клиент c таким номером не был найден, заказ со списанием баллов сохраняется в бд
// 			return err
// 		}
// 		return err
// 	}
// 	if userDB == userID {
// 		return errors.New("order already belong this user")
// 	}

// 	return errors.New("order already belong this user")

// }
