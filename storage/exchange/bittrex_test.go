package exchange

import (
	"reflect"
	"testing"
	"time"

	"github.com/nawa/cryptoexchange-dashboard/domain"
	"github.com/nawa/cryptoexchange-dashboard/utils"

	"github.com/Sirupsen/logrus"
	"github.com/h2non/gock"
	"github.com/nawa/cryptoexchange-dashboard/storage/exchange/testdata"
	assert "github.com/stretchr/testify/require"
	"github.com/toorop/go-bittrex"
)

const testAPIKey = "aaapppiiiKey"
const testAPISecret = "aaapppiiiSecret"

func TestNewBittrexExchange(t *testing.T) {
	exchange := NewBittrexExchange(testAPIKey, testAPISecret)
	assert.IsType(t, &bittrexExchange{}, exchange)
	assert.NotNil(t, exchange.(*bittrexExchange).bittrex)
	assert.NotNil(t, exchange.(*bittrexExchange).log)
}

func TestBittrexExchange_Ping(t *testing.T) {
	defer gock.Off()

	be := &bittrexExchange{
		bittrex: bittrex.New(testAPIKey, testAPISecret),
		log:     utils.NewDevNullLog(),
	}

	response := testdata.BittrexResponseSuccess([]bittrex.Balance{})

	gock.New("https://bittrex.com").
		//TODO check query params
		// AddMatcher(func(req *http.Request, gReq *gock.Request) (bool, error) {
		// 	apiKeyExists := req.URL.Query().Get("apikey") != ""
		// 	urlMatches := req.URL.Path == "/api/v1.1/account/getbalances"
		// 	fmt.Printf("%v - %v - %v\n", req.URL.Path, apiKeyExists, urlMatches)
		// 	return apiKeyExists && urlMatches, nil
		// }).
		Get("api/v1.1/account/getbalances").
		Reply(200).
		JSON(response)

	err := be.Ping()
	assert.NoError(t, err)
}

func TestBittrexExchange_GetBalance(t *testing.T) {
	type fields struct {
		bittrex *bittrex.Bittrex
		log     *logrus.Entry
	}
	tests := []struct {
		name    string
		fieldsF func() fields
		want    []domain.Balance
		wantErr bool
	}{
		{
			name: "correct",
			fieldsF: func() fields {
				response := testdata.BittrexResponseSuccess(testdata.BittrexMarketSummaries())

				gock.New("https://bittrex.com").
					Get("api/v1.1/public/getmarketsummaries").
					Reply(200).
					JSON(response)

				response = testdata.BittrexResponseSuccess(testdata.BittrexBalances())

				gock.New("https://bittrex.com").
					Get("api/v1.1/account/getbalances").
					Reply(200).
					JSON(response)

				return fields{
					bittrex: bittrex.New(testAPIKey, testAPISecret),
					log:     utils.NewDevNullLog(),
				}
			},
			want:    testdata.ModelBalances(),
			wantErr: false,
		},
		{
			name: "error in bittrex 'account|getbalances'",
			fieldsF: func() fields {
				response := testdata.BittrexResponseSuccess(testdata.BittrexMarketSummaries())

				gock.New("https://bittrex.com").
					Get("api/v1.1/public/getmarketsummaries").
					Reply(200).
					JSON(response)

				gock.New("https://bittrex.com").
					Get("api/v1.1/account/getbalances").
					Reply(200).
					JSON(testdata.BittrexResponseFailure)

				return fields{
					bittrex: bittrex.New(testAPIKey, testAPISecret),
					log:     utils.NewDevNullLog(),
				}
			},
			wantErr: true,
		},
		{
			name: "error in bittrex 'public|getmarketsummaries'",
			fieldsF: func() fields {
				gock.New("https://bittrex.com").
					Get("api/v1.1/public/getmarketsummaries").
					Reply(404)

				response := testdata.BittrexResponseSuccess(testdata.BittrexBalances())

				gock.New("https://bittrex.com").
					Get("api/v1.1/account/getbalances").
					Reply(200).
					JSON(response)

				return fields{
					bittrex: bittrex.New(testAPIKey, testAPISecret),
					log:     utils.NewDevNullLog(),
				}
			},
			wantErr: true,
		},
		{
			name: "BTC market is missing for currency",
			fieldsF: func() fields {
				response := testdata.BittrexResponseSuccess([]bittrex.MarketSummary{})
				gock.New("https://bittrex.com").
					Get("api/v1.1/public/getmarketsummaries").
					Reply(200).
					JSON(response)

				response = testdata.BittrexResponseSuccess([]bittrex.Balance{testdata.BittrexBalances()[1]})

				gock.New("https://bittrex.com").
					Get("api/v1.1/account/getbalances").
					Reply(200).
					JSON(response)

				return fields{
					bittrex: bittrex.New(testAPIKey, testAPISecret),
					log:     utils.NewDevNullLog(),
				}
			},
			wantErr: true,
		},
		{
			name: "USDT market is missing for BTC somehow",
			fieldsF: func() fields {
				response := testdata.BittrexResponseSuccess([]bittrex.MarketSummary{})
				gock.New("https://bittrex.com").
					Get("api/v1.1/public/getmarketsummaries").
					Reply(200).
					JSON(response)

				response = testdata.BittrexResponseSuccess([]bittrex.Balance{testdata.BittrexBalances()[0]})

				gock.New("https://bittrex.com").
					Get("api/v1.1/account/getbalances").
					Reply(200).
					JSON(response)

				return fields{
					bittrex: bittrex.New(testAPIKey, testAPISecret),
					log:     utils.NewDevNullLog(),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			fields := tt.fieldsF()
			be := &bittrexExchange{
				bittrex: fields.bittrex,
				log:     fields.log,
			}
			got, err := be.GetBalance()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("bittrexExchange.GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !modelBalancesEqual(got, tt.want) {
				t.Errorf("bittrexExchange.GetBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func modelBalancesEqual(got []domain.Balance, want []domain.Balance) bool {
	var zeroTime time.Time
	for i := range got {
		got[i].Time = zeroTime
	}
	return reflect.DeepEqual(got, want)
}

func TestBittrexExchange_GetMarketInfo(t *testing.T) {
	type fields struct {
		bittrex *bittrex.Bittrex
		log     *logrus.Entry
	}
	type args struct {
		market string
	}
	tests := []struct {
		name    string
		fieldsF func() fields
		args    args
		want    *domain.MarketInfo
		wantErr bool
	}{
		{
			name: "correct",
			fieldsF: func() fields {
				response := testdata.BittrexResponseSuccess(testdata.BittrexMarketSummaries())

				gock.New("https://bittrex.com").
					Get("api/v1.1/public/getmarketsummary").
					Reply(200).
					JSON(response)

				return fields{
					bittrex: bittrex.New(testAPIKey, testAPISecret),
					log:     utils.NewDevNullLog(),
				}
			},
			want:    testdata.ModelMarketInfo(),
			wantErr: false,
		},
		{
			name: "error in bittrex 'public|getmarketsummary'",
			fieldsF: func() fields {

				gock.New("https://bittrex.com").
					Get("api/v1.1/public/getmarketsummary").
					Reply(404)

				return fields{
					bittrex: bittrex.New(testAPIKey, testAPISecret),
					log:     utils.NewDevNullLog(),
				}
			},
			wantErr: true,
		},
		{
			name: "error when bittrex returns empty list of summaries",
			fieldsF: func() fields {
				response := testdata.BittrexResponseSuccess([]bittrex.MarketSummary{})

				gock.New("https://bittrex.com").
					Get("api/v1.1/public/getmarketsummary").
					Reply(200).
					JSON(response)

				return fields{
					bittrex: bittrex.New(testAPIKey, testAPISecret),
					log:     utils.NewDevNullLog(),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			fields := tt.fieldsF()
			be := &bittrexExchange{
				bittrex: fields.bittrex,
				log:     fields.log,
			}
			got, err := be.GetMarketInfo(tt.args.market)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("bittrexExchange.GetMarketInfo() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bittrexExchange.GetMarketInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBittrexExchange_GetOrders(t *testing.T) {
	type fields struct {
		bittrex *bittrex.Bittrex
		log     *logrus.Entry
	}
	tests := []struct {
		name    string
		fieldsF func() fields
		want    []domain.Order
		wantErr bool
	}{
		{
			name: "correct",
			fieldsF: func() fields {
				response := testdata.BittrexResponseSuccess(testdata.BittrexMarketSummaries())

				gock.New("https://bittrex.com").
					Get("api/v1.1/public/getmarketsummaries").
					Reply(200).
					JSON(response)

				response = testdata.BittrexResponseSuccess(testdata.BittrexOrders())

				gock.New("https://bittrex.com").
					Get("api/v1.1/account/getorderhistory").
					Reply(200).
					JSON(response)

				return fields{
					bittrex: bittrex.New(testAPIKey, testAPISecret),
					log:     utils.NewDevNullLog(),
				}
			},
			want:    testdata.ModelOrders(),
			wantErr: false,
		},
		{
			name: "correct if bittrex `account|getorderhistory` returns empty list of orders",
			fieldsF: func() fields {
				response := testdata.BittrexResponseSuccess(testdata.BittrexMarketSummaries())

				gock.New("https://bittrex.com").
					Get("api/v1.1/public/getmarketsummaries").
					Reply(200).
					JSON(response)

				response = testdata.BittrexResponseSuccess([]bittrex.Order{})

				gock.New("https://bittrex.com").
					Get("api/v1.1/account/getorderhistory").
					Reply(200).
					JSON(response)

				return fields{
					bittrex: bittrex.New(testAPIKey, testAPISecret),
					log:     utils.NewDevNullLog(),
				}
			},
			want:    []domain.Order{},
			wantErr: false,
		},
		{
			name: "error in bittrex 'public|getmarketsummaries'",
			fieldsF: func() fields {
				gock.New("https://bittrex.com").
					Get("api/v1.1/public/getmarketsummaries").
					Reply(404)

				response := testdata.BittrexResponseSuccess([]bittrex.Order{})

				gock.New("https://bittrex.com").
					Get("api/v1.1/account/getorderhistory").
					Reply(200).
					JSON(response)

				return fields{
					bittrex: bittrex.New(testAPIKey, testAPISecret),
					log:     utils.NewDevNullLog(),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			fields := tt.fieldsF()
			be := &bittrexExchange{
				bittrex: fields.bittrex,
				log:     fields.log,
			}
			got, err := be.GetOrders()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("bittrexExchange.GetOrders() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bittrexExchange.GetOrders() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
