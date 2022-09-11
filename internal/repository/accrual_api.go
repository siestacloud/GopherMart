package repository

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/siestacloud/gopherMart/internal/config"
	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/pkg"
	"github.com/sirupsen/logrus"
)

// OrderAccrual имплементирует работу с системой расчетов баллов лояльности
type AccrualAPI struct {
	cfg *config.Cfg
}

// NewOrderAccrual конструктор
func NewAccrualAPI(cfg *config.Cfg) *AccrualAPI {
	return &AccrualAPI{
		cfg: cfg,
	}
}

// GetOrderInfo взаимодействие с системой рассчета баллов лояльности
func (o *AccrualAPI) GetOrderAccrual(order *core.Order) error {
	var respErr core.AccrualAPIError

	client := resty.New().
		SetRetryCount(2).
		SetLogger(disableRestyDefLogger()).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(2 * time.Second)
	_, err := client.R().
		SetResult(order).
		ForceContentType("application/json").
		SetHeader("Content-Length", "0").
		SetError(&respErr).
		Get(o.cfg.AccrualSystemAddress + "/api/orders/" + order.Number)
	if err != nil {
		pkg.ErrPrint("repository", "err", err)

		return errors.New("unable GET accrual system API")
	}
	pkg.InfoPrint("repository", "ok", order)
	return nil
}

// disableRestyDefLogger отключаю дефолтное логирование пакета resty
func disableRestyDefLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return logger
}
