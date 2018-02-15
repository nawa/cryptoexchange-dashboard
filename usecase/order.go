package usecase

import (
	"github.com/Sirupsen/logrus"
	"github.com/nawa/cryptoexchange-wallet-info/model"
	"github.com/nawa/cryptoexchange-wallet-info/storage"
)

type OrderUsecase interface {
	GetActiveOrders() ([]model.Order, error)
}

type orderUsecase struct {
	exchange storage.Exchange
	log      *logrus.Entry
}

func NewOrderUsecase(exchange storage.Exchange) OrderUsecase {
	log := logrus.WithField("component", "orderUC")
	return &orderUsecase{
		exchange: exchange,
		log:      log,
	}
}

func (u *orderUsecase) GetActiveOrders() (orders []model.Order, err error) {
	orders, err = u.exchange.GetOrders()
	if err != nil {
		u.log.WithField("method", "GetActiveOrders").WithError(err).Error()
		return
	}

	return
}
