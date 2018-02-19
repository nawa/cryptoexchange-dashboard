package http

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iris-contrib/httpexpect"
	"github.com/kataras/iris/httptest"
	"github.com/pkg/errors"

	"github.com/nawa/cryptoexchange-dashboard/http/testdata"
	"github.com/nawa/cryptoexchange-dashboard/model"
	"github.com/nawa/cryptoexchange-dashboard/usecase/mocks"
)

type HTTPServerMock struct {
	Server     *Server
	BalanceUC  *mocks.MockBalanceUsecases
	OrderUC    *mocks.MockOrderUsecases
	HTTPExpect *httpexpect.Expect
}

func NewHTTPServerMock(t *testing.T, ctrl *gomock.Controller) *HTTPServerMock {
	balanceUC := mocks.NewMockBalanceUsecases(ctrl)
	orderUC := mocks.NewMockOrderUsecases(ctrl)
	server := NewServer(balanceUC, orderUC)

	return &HTTPServerMock{
		Server:     server,
		HTTPExpect: httptest.New(t, server.app),
		BalanceUC:  balanceUC,
		OrderUC:    orderUC,
	}
}

func TestBaseHandler_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	serverMock := NewHTTPServerMock(t, ctrl)

	response := serverMock.HTTPExpect.GET("/ping").Expect()

	response.Status(httptest.StatusOK)
	response.Body().Equal("pong")
}

func TestBalanceHandler_Hourly(t *testing.T) {
	runTestCases(t, []testCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchHourly("CUR1", 1).
					Return(testdata.CurrencyBalances()["CUR1"], nil)

				response := mock.HTTPExpect.GET("/balance/period/hourly/1").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusOK)
				response.Body().Equal(`{"CUR1":[{"amount":1,"btc":1,"usdt":2,"time":0},{"amount":2,"btc":2,"usdt":4,"time":3600}]}`)
			},
		}, {
			name: "correct with unknown currency",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchHourly("CUR1", 1).
					Return([]model.CurrencyBalance{}, nil)

				response := mock.HTTPExpect.GET("/balance/period/hourly/1").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusOK)
				response.Body().Equal(`{}`)
			},
		}, {
			name: "incorrect request: non int 'hours' path param",
			test: func(t *testing.T, mock *HTTPServerMock) {
				response := mock.HTTPExpect.GET("/balance/period/hourly/s").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusBadRequest)
			},
		}, {
			name: "incorrect request: 'hours' path param <=0",
			test: func(t *testing.T, mock *HTTPServerMock) {
				response := mock.HTTPExpect.GET("/balance/period/hourly/0").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusBadRequest)
			},
		}, {
			name: "incorrect request: 'currency' is missing",
			test: func(t *testing.T, mock *HTTPServerMock) {
				response := mock.HTTPExpect.GET("/balance/period/hourly/1").
					Expect()

				response.Status(httptest.StatusBadRequest)
			},
		},
		{
			name: "error from usecase level",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchHourly("CUR1", 1).
					Return(nil, errors.New("unexpected error"))

				response := mock.HTTPExpect.GET("/balance/period/hourly/1").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusInternalServerError)
			},
		},
	})
}

func TestBalanceHandler_Weekly(t *testing.T) {
	runTestCases(t, []testCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchWeekly("CUR1").
					Return(testdata.CurrencyBalances()["CUR1"], nil)

				response := mock.HTTPExpect.GET("/balance/period/weekly").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusOK)
				response.Body().Equal(`{"CUR1":[{"amount":1,"btc":1,"usdt":2,"time":0},{"amount":2,"btc":2,"usdt":4,"time":3600}]}`)
			},
		}, {
			name: "correct with unknown currency",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchWeekly("CUR1").
					Return([]model.CurrencyBalance{}, nil)

				response := mock.HTTPExpect.GET("/balance/period/weekly").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusOK)
				response.Body().Equal(`{}`)
			},
		}, {
			name: "incorrect request: 'currency' is missing",
			test: func(t *testing.T, mock *HTTPServerMock) {
				response := mock.HTTPExpect.GET("/balance/period/weekly").
					Expect()

				response.Status(httptest.StatusBadRequest)
			},
		}, {
			name: "error from usecase level",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchWeekly("CUR1").
					Return(nil, errors.New("unexpected error"))

				response := mock.HTTPExpect.GET("/balance/period/weekly").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusInternalServerError)
			},
		},
	})
}

func TestBalanceHandler_Monthly(t *testing.T) {
	runTestCases(t, []testCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchMonthly("CUR1").
					Return(testdata.CurrencyBalances()["CUR1"], nil)

				response := mock.HTTPExpect.GET("/balance/period/monthly").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusOK)
				response.Body().Equal(`{"CUR1":[{"amount":1,"btc":1,"usdt":2,"time":0},{"amount":2,"btc":2,"usdt":4,"time":3600}]}`)
			},
		}, {
			name: "correct with unknown currency",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchMonthly("CUR1").
					Return([]model.CurrencyBalance{}, nil)

				response := mock.HTTPExpect.GET("/balance/period/monthly").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusOK)
				response.Body().Equal(`{}`)
			},
		}, {
			name: "incorrect request: 'currency' is missing",
			test: func(t *testing.T, mock *HTTPServerMock) {
				response := mock.HTTPExpect.GET("/balance/period/monthly").
					Expect()

				response.Status(httptest.StatusBadRequest)
			},
		}, {
			name: "error from usecase level",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchMonthly("CUR1").
					Return(nil, errors.New("unexpected error"))

				response := mock.HTTPExpect.GET("/balance/period/monthly").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusInternalServerError)
			},
		},
	})
}

func TestBalanceHandler_All(t *testing.T) {
	runTestCases(t, []testCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchAll("CUR1").
					Return(testdata.CurrencyBalances()["CUR1"], nil)

				response := mock.HTTPExpect.GET("/balance/period/all").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusOK)
				response.Body().Equal(`{"CUR1":[{"amount":1,"btc":1,"usdt":2,"time":0},{"amount":2,"btc":2,"usdt":4,"time":3600}]}`)
			},
		}, {
			name: "correct with unknown currency",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchAll("CUR1").
					Return([]model.CurrencyBalance{}, nil)

				response := mock.HTTPExpect.GET("/balance/period/all").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusOK)
				response.Body().Equal(`{}`)
			},
		}, {
			name: "incorrect request: 'currency' is missing",
			test: func(t *testing.T, mock *HTTPServerMock) {
				response := mock.HTTPExpect.GET("/balance/period/all").
					Expect()

				response.Status(httptest.StatusBadRequest)
			},
		}, {
			name: "error from usecase level",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					FetchAll("CUR1").
					Return(nil, errors.New("unexpected error"))

				response := mock.HTTPExpect.GET("/balance/period/all").
					WithQuery("currency", "CUR1").
					Expect()

				response.Status(httptest.StatusInternalServerError)
			},
		},
	})
}

func TestBalanceHandler_ActiveCurrencies(t *testing.T) {
	runTestCases(t, []testCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					GetActiveCurrencies().
					Return(append(testdata.CurrencyBalances()["CUR1"], testdata.CurrencyBalances()["CUR2"]...), nil)

				response := mock.HTTPExpect.GET("/balance/active").
					Expect()

				response.Status(httptest.StatusOK)
				response.Body().Equal(`{"CUR1":[{"amount":1,"btc":1,"usdt":2,"time":0},{"amount":2,"btc":2,"usdt":4,"time":3600}],"CUR2":[{"amount":3,"btc":3,"usdt":5,"time":7200}]}`)
			},
		}, {
			name: "error from usecase level",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.BalanceUC.EXPECT().
					GetActiveCurrencies().
					Return(nil, errors.New("unexpected error"))

				response := mock.HTTPExpect.GET("/balance/active").
					Expect()

				response.Status(httptest.StatusInternalServerError)
			},
		},
	})
}

func TestOrderHandler_GetActiveOrders(t *testing.T) {
	runTestCases(t, []testCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.OrderUC.EXPECT().
					GetActiveOrders().
					Return(testdata.Orders(), nil)

				response := mock.HTTPExpect.GET("/order").
					Expect()

				response.Status(httptest.StatusOK)
				response.Body().Equal(`[{"market":"market1","market_link":"https://bittrex.com/Market/Index?MarketName=market1","time":0,"buy_rate":1,"amount":2,"sellnow_rate":3,"usdt_rate":4},{"market":"market2","market_link":"https://bittrex.com/Market/Index?MarketName=market2","time":3600,"buy_rate":5,"amount":6,"sellnow_rate":7,"usdt_rate":8}]`)
			},
		}, {
			name: "error from usecase level",
			test: func(t *testing.T, mock *HTTPServerMock) {
				mock.OrderUC.EXPECT().
					GetActiveOrders().
					Return(nil, errors.New("unexpected error"))

				response := mock.HTTPExpect.GET("/order").
					Expect()

				response.Status(httptest.StatusInternalServerError)
			},
		},
	})
}

type testCase struct {
	name string
	test func(t *testing.T, mock *HTTPServerMock)
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewHTTPServerMock(t, ctrl)
			tc.test(t, mock)
		})
	}
}
