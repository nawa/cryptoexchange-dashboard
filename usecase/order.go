package usecase

import (
	"github.com/Sirupsen/logrus"
	"github.com/nawa/cryptoexchange-dashboard/domain"
	"github.com/nawa/cryptoexchange-dashboard/storage"
)

type OrderUsecases interface {
	GetActiveOrders() ([]domain.Order, error)
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

func (u *orderUsecases) GetActiveOrders() ([]domain.Order, error) {
	orders, err := u.exchange.GetOrders()
	if err != nil {
		u.log.WithField("method", "GetActiveOrders").WithError(err).Error()
		return nil, err
	}

	return orders, nil
}
