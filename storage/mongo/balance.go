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

func (s *balanceStorage) Save(balances ...storage.Balance) error {
	db, closeSession := s.getDB()
	defer closeSession()

	convert := make([]interface{}, len(balances))
	for i := range balances {
		convert[i] = balances[i]
	}
	return db.C("balance").Insert(convert...)
}

func (s *balanceStorage) FetchHourly(currency string, hours int) (balance []storage.Balance, err error) {
	db, closeSession := s.getDB()
	defer closeSession()

	period := time.Now().Add(-1 * time.Hour * time.Duration(hours))

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

func (s *balanceStorage) FetchWeekly(currency string) (balance []storage.Balance, err error) {
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
		bson.M{
			"$match": bson.M{
				"time": bson.M{
					"$gte": period,
				},
				"currency": currency,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"year":  bson.M{"$year": "$time"},
					"month": bson.M{"$month": "$time"},
					"day":   bson.M{"$dayOfMonth": "$time"},
					"hour":  bson.M{"$hour": "$time"},
					"each_5min": bson.M{
						"$subtract": []bson.M{
							bson.M{"$minute": "$time"},
							bson.M{"$mod": []interface{}{bson.M{"$minute": "$time"}, 5}},
						},
					},
				},
				"time":        bson.M{"$first": "$time"},
				"amount":      bson.M{"$first": "$amount"},
				"btc_amount":  bson.M{"$first": "$btc_amount"},
				"usdt_amount": bson.M{"$first": "$usdt_amount"},
			},
		},
		bson.M{"$sort": bson.M{"time": -1}},
	})

	err = pipe.
		All(&balance)

	return
}

func (s *balanceStorage) FetchMonthly(currency string) (balance []storage.Balance, err error) {
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
		bson.M{
			"$match": bson.M{
				"time": bson.M{
					"$gte": period,
				},
				"currency": currency,
			},
		},
		bson.M{
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
		bson.M{"$sort": bson.M{"time": -1}},
	})

	err = pipe.
		All(&balance)

	return
}

func (s *balanceStorage) FetchAll(currency string) (balance []storage.Balance, err error) {
	panic("not implemented")
}

func (s *balanceStorage) GetActiveCurrencies() (balance []storage.Balance, err error) {
	//TODO make in one call to mongo
	db, closeSession := s.getDB()
	defer closeSession()

	type lastTime struct {
		Time time.Time `bson:"time"`
	}
	var t []lastTime

	err = db.C("balance").
		Find(bson.M{}).
		Sort("-time").
		Limit(1).
		All(&t)

	if err != nil {
		return
	}

	if len(t) == 0 {
		return
	}

	q := bson.M{
		"time": t[0].Time,
	}

	err = db.C("balance").
		Find(q).
		All(&balance)

	return
}
