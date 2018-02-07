package mongo

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/nawa/cryptoexchange-wallet-info/storage"
)

type balanceStorage struct {
	baseStorage
}

func NewBalanceStorage(session *mgo.Session, refreshSession bool) storage.BalanceStorage {
	return &balanceStorage{
		baseStorage{
			baseSession:    session,
			refreshSession: refreshSession,
		},
	}
}

func (s *balanceStorage) Save(balances ...storage.Balance) error {
	db, closeSession := s.getDB()
	defer closeSession()

	convert := make([]interface{}, len(balances))
	for i := range balances {
		convert[i] = balances[i]
	}
	return db.C("balance").Insert(convert...)
}

func (s *balanceStorage) FetchDaily(currency string) (balance []storage.Balance, err error) {
	db, closeSession := s.getDB()
	defer closeSession()

	period := time.Now().Add(-1 * time.Hour * 24)

	q := bson.M{
		"time": bson.M{
			"$gte": period,
		},
		"currency": currency,
	}
	err = db.C("balance").
		Find(q).
		Sort("-time").
		All(&balance)
	return
}
