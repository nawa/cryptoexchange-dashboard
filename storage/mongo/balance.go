package mongo

import (
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

func (s *balanceStorage) Find() (balance []storage.Balance, err error) {
	db, closeSession := s.getDB()
	defer closeSession()

	q := bson.M{}
	err = db.C("balance").Find(q).All(&balance)
	return
}
