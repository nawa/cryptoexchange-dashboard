package exchange

import "github.com/nawa/cryptoexchange-wallet-info/sync/model"

type Exchange interface {
	GetBalance() (*model.Balance, error)
	Ping() error
}
