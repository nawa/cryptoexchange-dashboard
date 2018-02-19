package usecase

import (
	"github.com/Sirupsen/logrus"
	"github.com/nawa/cryptoexchange-dashboard/model"
	"github.com/nawa/cryptoexchange-dashboard/storage"
)

type OrderUsecases interface {
	GetActiveOrders() ([]model.Order, error)
}

type orderUsecases struct {
	exchange storage.Exchange
	log      *logrus.Entry
}

func NewOrderUsecase(exchange storage.Exchange) OrderUsecases {
	log := logrus.WithField("component", "orderUC")
	return &orderUsecases{
		exchange: exchange,
		log:      log,
	}
}

func (u *orderUsecases) GetActiveOrders() (orders []model.Order, err error) {
	orders, err = u.exchange.GetOrders()
	if err != nil {
		u.log.WithField("method", "GetActiveOrders").WithError(err).Error()
		return
	}

	return
}
