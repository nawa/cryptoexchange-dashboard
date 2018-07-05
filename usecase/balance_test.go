package usecase

import (
	"reflect"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	assert "github.com/stretchr/testify/require"

	"github.com/nawa/cryptoexchange-dashboard/domain"
	"github.com/nawa/cryptoexchange-dashboard/storage"
	"github.com/nawa/cryptoexchange-dashboard/storage/mocks"
	"github.com/nawa/cryptoexchange-dashboard/usecase/testdata"
	"github.com/nawa/cryptoexchange-dashboard/utils"
)

var errExpected = errors.New("some_test_error")

func TestBalanceUsecases_StartSyncFromExchangePeriodically(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	balanceStorage := mocks.NewMockBalanceStorage(ctrl)
	exchange := mocks.NewMockExchange(ctrl)

	exchange.EXPECT().GetBalance().
		Return(testdata.Balances(), nil).
		MinTimes(10).
		MaxTimes(20)

	balanceStorage.EXPECT().
		Save(testdata.BalancesWithTotal()).
		Return(nil).
		MinTimes(10).
		MaxTimes(20)

	balanceUC := NewBalanceUsecase(exchange, balanceStorage)

	stop, err := balanceUC.StartSyncFromExchangePeriodically(time.Millisecond * 10)
	assert.NoError(t, err)

	// <= 20 operations
	time.Sleep(time.Millisecond * 200)

	stop()

	//wait a little after stop
	time.Sleep(time.Millisecond * 100)
}

func TestBalanceUsecases_SyncFromExchange(t *testing.T) {
	type fields struct {
		exchange       storage.Exchange
		balanceStorage storage.BalanceStorage
		log            *logrus.Entry
	}
	tests := []struct {
		name    string
		fieldsF func(ctrl *gomock.Controller) fields
		wantErr bool
	}{
		{
			name: "correct",
			fieldsF: func(ctrl *gomock.Controller) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				exchange.EXPECT().GetBalance().
					Return(testdata.Balances(), nil).
					Times(1)

				balanceStorage.EXPECT().
					Save(testdata.BalancesWithTotal()).
					Return(nil).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			wantErr: false,
		},
		{
			name: "error in exchange",
			fieldsF: func(ctrl *gomock.Controller) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				exchange.EXPECT().GetBalance().
					Return(nil, errExpected).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			wantErr: true,
		},
		{
			name: "error in balanceStorage",
			fieldsF: func(ctrl *gomock.Controller) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				exchange.EXPECT().GetBalance().
					Return(testdata.Balances(), nil).
					Times(1)

				balanceStorage.EXPECT().
					Save(testdata.BalancesWithTotal()).
					Return(errExpected).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			wantErr: true,
		},
		{
			name: "log with debug level",
			fieldsF: func(ctrl *gomock.Controller) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				exchange.EXPECT().GetBalance().
					Return(testdata.Balances(), nil).
					Times(1)

				balanceStorage.EXPECT().
					Save(testdata.BalancesWithTotal()).
					Return(nil).
					Times(1)

				logger := utils.NewDevNullLog()

				logger.Level = logrus.DebugLevel

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            logger,
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := tt.fieldsF(ctrl)
			u := &balanceUsecases{
				exchange:       fields.exchange,
				balanceStorage: fields.balanceStorage,
				log:            fields.log,
			}
			if err := u.SyncFromExchange(); err != nil {
				if !tt.wantErr {
					t.Errorf("balanceUsecases.SyncFromExchange() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err != errExpected {
					t.Errorf("balanceUsecases.SyncFromExchange() error = %v, expected error %v", err, errExpected)
				}
			}
		})
	}
}

func TestBalanceUsecases_FetchHourly(t *testing.T) {
	type fields struct {
		exchange       storage.Exchange
		balanceStorage storage.BalanceStorage
		log            *logrus.Entry
	}
	type args struct {
		currency string
		hours    int
	}
	tests := []struct {
		name         string
		fieldsF      func(ctrl *gomock.Controller, args args) fields
		args         args
		wantBalances []domain.Balance
		wantErr      bool
	}{
		{
			name: "correct",
			fieldsF: func(ctrl *gomock.Controller, args args) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				balanceStorage.EXPECT().
					FetchHourly(args.currency, args.hours).
					Return(testdata.Balances(), nil).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			args: args{
				currency: "CURS",
				hours:    2,
			},
			wantBalances: testdata.Balances(),
			wantErr:      false,
		},
		{
			name: "error in balanceStorage",
			fieldsF: func(ctrl *gomock.Controller, args args) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				balanceStorage.EXPECT().
					FetchHourly(args.currency, args.hours).
					Return(nil, errExpected).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			args: args{
				currency: "CURS",
				hours:    2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := tt.fieldsF(ctrl, tt.args)
			u := &balanceUsecases{
				exchange:       fields.exchange,
				balanceStorage: fields.balanceStorage,
				log:            fields.log,
			}
			gotBalances, err := u.FetchHourly(tt.args.currency, tt.args.hours)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("balanceUsecases.FetchHourly() error = %v, wantErr %v", err, tt.wantErr)
				} else if err != errExpected {
					t.Errorf("balanceUsecases.FetchHourly() error = %v, expected error %v", err, errExpected)
				}
				return
			}
			if !reflect.DeepEqual(gotBalances, tt.wantBalances) {
				t.Errorf("balanceUsecases.FetchHourly() = %v, want %v", gotBalances, tt.wantBalances)
			}
		})
	}
}

func TestBalanceUsecases_FetchWeekly(t *testing.T) {
	type fields struct {
		exchange       storage.Exchange
		balanceStorage storage.BalanceStorage
		log            *logrus.Entry
	}
	type args struct {
		currency string
	}
	tests := []struct {
		name         string
		fieldsF      func(ctrl *gomock.Controller, args args) fields
		args         args
		wantBalances []domain.Balance
		wantErr      bool
	}{
		{
			name: "correct",
			fieldsF: func(ctrl *gomock.Controller, args args) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				balanceStorage.EXPECT().
					FetchWeekly(args.currency).
					Return(testdata.Balances(), nil).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			args:         args{currency: "CURS"},
			wantBalances: testdata.Balances(),
			wantErr:      false,
		},
		{
			name: "error in balanceStorage",
			fieldsF: func(ctrl *gomock.Controller, args args) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				balanceStorage.EXPECT().
					FetchWeekly(args.currency).
					Return(nil, errExpected).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			args:    args{currency: "CURS"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := tt.fieldsF(ctrl, tt.args)
			u := &balanceUsecases{
				exchange:       fields.exchange,
				balanceStorage: fields.balanceStorage,
				log:            fields.log,
			}
			gotBalances, err := u.FetchWeekly(tt.args.currency)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("balanceUsecases.FetchWeekly() error = %v, wantErr %v", err, tt.wantErr)
				} else if err != errExpected {
					t.Errorf("balanceUsecases.FetchWeekly() error = %v, expected error %v", err, errExpected)
				}
				return
			}
			if !reflect.DeepEqual(gotBalances, tt.wantBalances) {
				t.Errorf("balanceUsecases.FetchWeekly() = %v, want %v", gotBalances, tt.wantBalances)
			}
		})
	}
}

func TestBalanceUsecases_FetchMonthly(t *testing.T) {
	type fields struct {
		exchange       storage.Exchange
		balanceStorage storage.BalanceStorage
		log            *logrus.Entry
	}
	type args struct {
		currency string
	}
	tests := []struct {
		name         string
		fieldsF      func(ctrl *gomock.Controller, args args) fields
		args         args
		wantBalances []domain.Balance
		wantErr      bool
	}{
		{
			name: "correct",
			fieldsF: func(ctrl *gomock.Controller, args args) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				balanceStorage.EXPECT().
					FetchMonthly(args.currency).
					Return(testdata.Balances(), nil).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			args:         args{currency: "CURS"},
			wantBalances: testdata.Balances(),
			wantErr:      false,
		},
		{
			name: "error in balanceStorage",
			fieldsF: func(ctrl *gomock.Controller, args args) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				balanceStorage.EXPECT().
					FetchMonthly(args.currency).
					Return(nil, errExpected).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			args:    args{currency: "CURS"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := tt.fieldsF(ctrl, tt.args)
			u := &balanceUsecases{
				exchange:       fields.exchange,
				balanceStorage: fields.balanceStorage,
				log:            fields.log,
			}
			gotBalances, err := u.FetchMonthly(tt.args.currency)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("balanceUsecases.FetchMonthly() error = %v, wantErr %v", err, tt.wantErr)
				} else if err != errExpected {
					t.Errorf("balanceUsecases.FetchMonthly() error = %v, expected error %v", err, errExpected)
				}
				return
			}
			if !reflect.DeepEqual(gotBalances, tt.wantBalances) {
				t.Errorf("balanceUsecases.FetchMonthly() = %v, want %v", gotBalances, tt.wantBalances)
			}
		})
	}
}

func TestBalanceUsecases_FetchAll(t *testing.T) {
	type fields struct {
		exchange       storage.Exchange
		balanceStorage storage.BalanceStorage
		log            *logrus.Entry
	}
	type args struct {
		currency string
	}
	tests := []struct {
		name         string
		fieldsF      func(ctrl *gomock.Controller, args args) fields
		args         args
		wantBalances []domain.Balance
		wantErr      bool
	}{
		{
			name: "correct",
			fieldsF: func(ctrl *gomock.Controller, args args) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				balanceStorage.EXPECT().
					FetchAll(args.currency).
					Return(testdata.Balances(), nil).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			args:         args{currency: "CURS"},
			wantBalances: testdata.Balances(),
			wantErr:      false,
		},
		{
			name: "error in balanceStorage",
			fieldsF: func(ctrl *gomock.Controller, args args) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				balanceStorage.EXPECT().
					FetchAll(args.currency).
					Return(nil, errExpected).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			args:    args{currency: "CURS"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := tt.fieldsF(ctrl, tt.args)
			u := &balanceUsecases{
				exchange:       fields.exchange,
				balanceStorage: fields.balanceStorage,
				log:            fields.log,
			}
			gotBalances, err := u.FetchAll(tt.args.currency)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("balanceUsecases.FetchAll() error = %v, wantErr %v", err, tt.wantErr)
				} else if err != errExpected {
					t.Errorf("balanceUsecases.FetchAll() error = %v, expected error %v", err, errExpected)
				}
				return
			}
			if !reflect.DeepEqual(gotBalances, tt.wantBalances) {
				t.Errorf("balanceUsecases.FetchAll() = %v, want %v", gotBalances, tt.wantBalances)
			}
		})
	}
}

func TestBalanceUsecases_GetActiveCurrencies(t *testing.T) {
	type fields struct {
		exchange       storage.Exchange
		balanceStorage storage.BalanceStorage
		log            *logrus.Entry
	}
	tests := []struct {
		name         string
		fieldsF      func(ctrl *gomock.Controller) fields
		wantBalances []domain.Balance
		wantErr      bool
	}{
		{
			name: "correct",
			fieldsF: func(ctrl *gomock.Controller) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				balanceStorage.EXPECT().
					GetActiveCurrencies().
					Return(testdata.Balances(), nil).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			wantBalances: testdata.Balances(),
			wantErr:      false,
		}, {
			name: "error in storageBalance",
			fieldsF: func(ctrl *gomock.Controller) fields {
				balanceStorage := mocks.NewMockBalanceStorage(ctrl)
				exchange := mocks.NewMockExchange(ctrl)

				balanceStorage.EXPECT().
					GetActiveCurrencies().
					Return(nil, errExpected).
					Times(1)

				return fields{
					exchange:       exchange,
					balanceStorage: balanceStorage,
					log:            utils.NewDevNullLog(),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := tt.fieldsF(ctrl)
			u := &balanceUsecases{
				exchange:       fields.exchange,
				balanceStorage: fields.balanceStorage,
				log:            fields.log,
			}
			gotBalances, err := u.GetActiveCurrencies()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("balanceUsecases.GetActiveCurrencies() error = %v, wantErr %v", err, tt.wantErr)
				} else if err != errExpected {
					t.Errorf("balanceUsecases.GetActiveCurrencies() error = %v, expected error %v", err, errExpected)
				}
				return
			}
			if !reflect.DeepEqual(gotBalances, tt.wantBalances) {
				t.Errorf("balanceUsecases.GetActiveCurrencies() = %v, want %v", gotBalances, tt.wantBalances)
			}
		})
	}
}
