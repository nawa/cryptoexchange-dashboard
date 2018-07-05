package mongo

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"

	"github.com/nawa/cryptoexchange-dashboard/domain"
	"github.com/nawa/cryptoexchange-dashboard/storage"
)

type balanceStorage struct {
	baseStorage
}

type balance struct {
	Exchange   string    `bson:"exchange"`
	Currency   string    `bson:"currency"`
	Amount     float64   `bson:"amount"`
	BTCAmount  float64   `bson:"btc_amount"`
	USDTAmount float64   `bson:"usdt_amount"`
	Time       time.Time `bson:"time"`
}

func NewBalanceStorage(session *mgo.Session, refreshSession bool) storage.BalanceStorage {
	return &balanceStorage{
		baseStorage{
			baseSession:    session,
			refreshSession: refreshSession,
		},
	}
}

func (s *balanceStorage) Init() (err error) {
	db, closeSession := s.getDB()
	defer closeSession()

	c := db.C("balance")
	err = c.EnsureIndex(mgo.Index{
		Name:       "time_curr_idx",
		Key:        []string{"-time", "currency"},
		Unique:     false,
		Background: true,
	})

	if err != nil {
		return
	}

	err = c.EnsureIndex(mgo.Index{
		Name:       "time_idx",
		Key:        []string{"-time"},
		Unique:     false,
		Background: true,
	})

	if err != nil {
		return
	}

	err = c.EnsureIndex(mgo.Index{
		Name:       "curr_idx",
		Key:        []string{"curr"},
		Unique:     false,
		Background: true,
	})

	return
}

func (s *balanceStorage) Save(balance ...domain.Balance) error {
	db, closeSession := s.getDB()
	defer closeSession()

	return db.C("balance").Insert(convertBalancesFromModel(balance...)...)
}

func (s *balanceStorage) FetchHourly(currency string, hours int) ([]domain.Balance, error) {
	db, closeSession := s.getDB()
	defer closeSession()

	period := time.Now().Add(-1 * time.Hour * time.Duration(hours))

	q := bson.M{
		"time": bson.M{
			"$gte": period,
		},
		"currency": currency,
	}
	var balances []balance

	err := db.C("balance").
		Find(q).
		Sort("-time").
		All(&balances)
	if err != nil {
		return nil, err
	}

	return convertBalancesToModel(balances...), nil
}

func (s *balanceStorage) FetchWeekly(currency string) ([]domain.Balance, error) {
	db, closeSession := s.getDB()
	defer closeSession()

	// 	db.balance.aggregate(
	//     {
	//         $match: {
	//             "time": {
	//                 $gte: new Date((new Date().getTime() - (7 * 24 * 60 * 60 * 1000)))
	//             },
	//             "currency": "total"
	//         }
	//     },
	//     {
	//         $group: {
	//             "_id" : {
	//                 "year": {"$year":"$time"}, "month": {"$month":"$time"}, "day": {"$dayOfMonth":"$time"}, "hour": {"$hour":"$time"},
	//                 "each_5min": {
	//                     "$subtract": [
	//                         { "$minute": "$time" },
	//                         { "$mod": [{ "$minute": "$time"}, 5] }
	//                     ]
	//                 }
	//             },
	//             "total": {$sum: 1},
	//             "time": {$first: "$time"},
	//             "amount": {$first: "$amount"},
	//             "btc_amount": {$first: "$btc_amount"},
	//             "usdt_amount": {$first: "$usdt_amount"}
	//         }
	//     },
	//     { $sort : { time : -1}}
	// )

	period := time.Now().Add(-1 * time.Hour * 24 * 7)

	pipe := db.C("balance").Pipe([]bson.M{
		{
			"$match": bson.M{
				"time": bson.M{
					"$gte": period,
				},
				"currency": currency,
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"year":  bson.M{"$year": "$time"},
					"month": bson.M{"$month": "$time"},
					"day":   bson.M{"$dayOfMonth": "$time"},
					"hour":  bson.M{"$hour": "$time"},
					"each_5min": bson.M{
						"$subtract": []bson.M{
							{"$minute": "$time"},
							{"$mod": []interface{}{bson.M{"$minute": "$time"}, 5}},
						},
					},
				},
				"time":        bson.M{"$first": "$time"},
				"amount":      bson.M{"$first": "$amount"},
				"btc_amount":  bson.M{"$first": "$btc_amount"},
				"usdt_amount": bson.M{"$first": "$usdt_amount"},
			},
		},
		{"$sort": bson.M{"time": -1}},
	})

	var balances []balance

	err := pipe.
		All(&balances)

	if err != nil {
		return nil, err
	}

	return convertBalancesToModel(balances...), nil
}

func (s *balanceStorage) FetchMonthly(currency string) ([]domain.Balance, error) {
	db, closeSession := s.getDB()
	defer closeSession()

	// 	db.balance.aggregate(
	//     {
	//         $match: {
	//             "time": {
	//                 $gte: new Date((new Date().getTime() - (30 * 24 * 60 * 60 * 1000)))
	//             },
	//             "currency": "total"
	//         }
	//     },
	//     {
	//         $group: {
	//             "_id" : {
	//                 "year": {"$year":"$time"}, "month": {"$month":"$time"}, "day": {"$dayOfMonth":"$time"}, "hour": {"$hour":"$time"}
	//             },
	//             "total": {$sum: 1},
	//             "time": {$first: "$time"},
	//             "amount": {$first: "$amount"},
	//             "btc_amount": {$first: "$btc_amount"},
	//             "usdt_amount": {$first: "$usdt_amount"}
	//         }
	//     },
	//     { $sort : { time : -1}}
	// )

	period := time.Now().Add(-1 * time.Hour * 24 * 30)

	pipe := db.C("balance").Pipe([]bson.M{
		{
			"$match": bson.M{
				"time": bson.M{
					"$gte": period,
				},
				"currency": currency,
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"year":  bson.M{"$year": "$time"},
					"month": bson.M{"$month": "$time"},
					"day":   bson.M{"$dayOfMonth": "$time"},
					"hour":  bson.M{"$hour": "$time"},
				},
				"time":        bson.M{"$first": "$time"},
				"amount":      bson.M{"$first": "$amount"},
				"btc_amount":  bson.M{"$first": "$btc_amount"},
				"usdt_amount": bson.M{"$first": "$usdt_amount"},
			},
		},
		{"$sort": bson.M{"time": -1}},
	})

	var balances []balance

	err := pipe.
		All(&balances)

	if err != nil {
		return nil, err
	}

	return convertBalancesToModel(balances...), nil
}

func (s *balanceStorage) FetchAll(currency string) ([]domain.Balance, error) {
	panic("not implemented")
}

func (s *balanceStorage) GetActiveCurrencies() ([]domain.Balance, error) {
	//TODO make in one call to mongo
	db, closeSession := s.getDB()
	defer closeSession()

	type lastTime struct {
		Time time.Time `bson:"time"`
	}
	var t []lastTime

	err := db.C("balance").
		Find(bson.M{}).
		Sort("-time").
		Limit(1).
		All(&t)

	if err != nil {
		return nil, err
	}

	if len(t) == 0 {
		return nil, errors.New("last time not found for currency")
	}

	q := bson.M{
		"time": t[0].Time,
	}

	var balances []balance

	err = db.C("balance").
		Find(q).
		All(&balances)

	if err != nil {
		return nil, err
	}

	return convertBalancesToModel(balances...), nil
}

func convertBalancesFromModel(balances ...domain.Balance) (result []interface{}) {
	for _, b := range balances {
		result = append(result, balance{
			Exchange:   string(b.Exchange),
			Currency:   b.Currency,
			Amount:     b.Amount,
			BTCAmount:  b.BTCAmount,
			USDTAmount: b.USDTAmount,
			Time:       b.Time,
		})
	}
	return result
}

func convertBalancesToModel(balances ...balance) (result []domain.Balance) {
	for _, b := range balances {
		result = append(result, domain.Balance{
			Exchange:   domain.ExchangeType(b.Exchange),
			Currency:   b.Currency,
			Amount:     b.Amount,
			BTCAmount:  b.BTCAmount,
			USDTAmount: b.USDTAmount,
			Time:       b.Time,
		})
	}
	return result
}
