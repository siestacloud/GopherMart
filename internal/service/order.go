package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/siestacloud/gopherMart/internal/config"
	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/internal/repository"
	"github.com/siestacloud/gopherMart/pkg"
	"github.com/sirupsen/logrus"
)

const (
	statusNEW        string = "NEW"
	statusINVALID    string = "INVALID"
	statusPROCESSED  string = "PROCESSED"
	statusPROCESSING string = "PROCESSING"
)

// OrderService реализация бизнес логики обработки номера заказа
type OrderService struct {
	cfg *config.Cfg

	repo repository.Order
}

//NewOrderService конструктор
func NewOrderService(cfg *config.Cfg, repo repository.Order) *OrderService {
	return &OrderService{
		cfg:  cfg,
		repo: repo,
	}
}

//Create проверка номера заказа по алгоритму ЛУНА и сохранение его в базе (с привязкой к конкретному пользователю)
func (o *OrderService) Create(userID int, order core.Order) error {

	// * проверка номера заказа по алгоритму Луна
	if err := pkg.Valid(order.ID); err != nil {
		pkg.WarnPrint("service", "lune alg err", err)
		return err
	}

	userDB, err := o.repo.GetUserByOrder(order.ID) // * попытка определить клиента по номеру заказа
	if err != nil {
		pkg.WarnPrint("service", "get user by order", err)

		currentTime := time.Now().Format(time.RFC3339)

		url(o.cfg.AccrualSystemAddress, order.ID)
		if err = o.repo.Create(userID, order, statusPROCESSED, currentTime); err != nil { // * клиент c таким номером не был найден, заказ сохраняется в бд
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

type APIError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func url(address, orderID string) {
	var responseErr APIError
	pkg.InfoPrint("url", "http://"+address+"/orders/"+orderID)
	logger := logrus.New()
	logger.Out = ioutil.Discard
	client := resty.New().SetRetryCount(2).SetLogger(logger).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(2 * time.Second)
	rec, err := client.R().
		SetError(&responseErr).SetDoNotParseResponse(false).
		SetBody(nil).
		Post("http://" + address + "/api/orders/" + orderID)
	if err != nil {
		// fmt.Println("resp err:  ", responseErr)
		log.Println("AGENT resp err:: ", err)
	}

	pkg.InfoPrint("url", string(rec.Body()))
}
