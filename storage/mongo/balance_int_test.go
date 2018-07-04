// +build integration_test

package mongo_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/nawa/cryptoexchange-dashboard/domain"
	assert "github.com/stretchr/testify/require"

	"github.com/nawa/cryptoexchange-dashboard/storage"
	"github.com/nawa/cryptoexchange-dashboard/storage/mongo"
)

const EnvDbTestURL = "DB_TEST_URL"

var (
	session        *mgo.Session
	balanceStorage storage.BalanceStorage
)

func TestMain(m *testing.M) {
	dbURL, ok := os.LookupEnv(EnvDbTestURL)
	if !ok {
		log.Fatalf("environment variable '%s' is required", EnvDbTestURL)
	}

	dialInfo, err := mgo.ParseURL(dbURL)
	dialInfo.Timeout = time.Second * 10
	if err != nil {
		log.Fatalf("URL to mongo instance is incorrect: %s", err.Error())
	}

	session, err = mgo.DialWithInfo(dialInfo)
	defer session.Close()
	if err != nil {
		log.Fatalf("can't connect to mongo instance: %s", err.Error())
	}

	balanceStorage = mongo.NewBalanceStorage(session, true)
	err = balanceStorage.Init()
	if err != nil {
		log.Fatalf("can't instantiate balanceStorage: %s", err.Error())
	}

	err = cleanupData(session)
	if err != nil {
		log.Fatalf("can't cleanup data before tests: %s", err.Error())
	}

	code := m.Run()

	session.DB("").
		C("balance").
		DropCollection()

	os.Exit(code)
}

func TestBalanceStorage_FetchAll(t *testing.T) {
	defer func() {
		err := recover()
		assert.Equal(t, "not implemented", err)
	}()
	balanceStorage.FetchAll("USDT")
}

func TestBalanceStorage_FetchHourly(t *testing.T) {
	assert.NoError(t, cleanupData(session))

	now := time.Now()
	balances := []*domain.Balance{
		{
			Currencies: []domain.CurrencyBalance{
				{
					Currency:   "CUR1",
					Amount:     1,
					BTCAmount:  2,
					USDTAmount: 3,
					Time:       now.Add(-5 * time.Hour),
				},
				{
					Currency:   "CUR2",
					Amount:     4,
					BTCAmount:  5,
					USDTAmount: 6,
					Time:       now,
				},
			},
			Exchange:   domain.ExchangeTypeBittrex,
			BTCAmount:  100,
			USDTAmount: 1000,
			Time:       now,
		},
		{
			Currencies: []domain.CurrencyBalance{
				{
					Currency:   "CUR1",
					Amount:     1,
					BTCAmount:  2,
					USDTAmount: 3,
					Time:       now,
				},
				{
					Currency:   "CUR2",
					Amount:     4,
					BTCAmount:  5,
					USDTAmount: 6,
					Time:       now.Add(-5 * time.Hour),
				},
			},
			Exchange:   domain.ExchangeTypeBittrex,
			BTCAmount:  100,
			USDTAmount: 1000,
			Time:       now,
		},
	}

	for _, balance := range balances {
		err := balanceStorage.Save(balance)
		assert.NoError(t, err)
	}

	storageBalances, err := balanceStorage.FetchHourly("total", 2)
	assert.NoError(t, err)
	assert.Len(t, storageBalances, 2)

	storageBalances, err = balanceStorage.FetchHourly("CUR1", 2)
	assert.NoError(t, err)
	assert.Len(t, storageBalances, 1)

	storageBalances, err = balanceStorage.FetchHourly("CUR2", 2)
	assert.NoError(t, err)
	assert.Len(t, storageBalances, 1)
}

func cleanupData(session *mgo.Session) (err error) {
	_, err = session.DB("").
		C("balance").
		RemoveAll(bson.M{})

	return
}
