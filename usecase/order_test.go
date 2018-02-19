package usecase

import (
	"reflect"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/nawa/cryptoexchange-dashboard/model"
	"github.com/nawa/cryptoexchange-dashboard/storage"
	"github.com/nawa/cryptoexchange-dashboard/storage/mocks"
	"github.com/nawa/cryptoexchange-dashboard/usecase/testdata"
	"github.com/nawa/cryptoexchange-dashboard/utils"
	assert "github.com/stretchr/testify/require"
)

func TestOrderUsecases_GetActiveOrders(t *testing.T) {
	type fields struct {
		exchange storage.Exchange
		log      *logrus.Entry
	}
	tests := []struct {
		name       string
		fieldsF    func(ctrl *gomock.Controller) fields
		wantOrders []model.Order
		wantErr    bool
	}{
		{
			name: "correct",
			fieldsF: func(ctrl *gomock.Controller) fields {
				exchange := mocks.NewMockExchange(ctrl)

				exchange.EXPECT().
					GetOrders().
					Return(testdata.Orders(), nil).
					Times(1)

				return fields{
					exchange: exchange,
					log:      utils.NewDevNullLog(),
				}
			},
			wantOrders: testdata.Orders(),
			wantErr:    false,
		}, {
			name: "error in exchange",
			fieldsF: func(ctrl *gomock.Controller) fields {
				exchange := mocks.NewMockExchange(ctrl)

				exchange.EXPECT().
					GetOrders().
					Return(nil, errExpected).
					Times(1)

				return fields{
					exchange: exchange,
					log:      utils.NewDevNullLog(),
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
			u := &orderUsecases{
				exchange: fields.exchange,
				log:      fields.log,
			}
			gotOrders, err := u.GetActiveOrders()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("orderUsecases.GetActiveOrders() error = %v, wantErr %v", err, tt.wantErr)
				} else if err != errExpected {
					t.Errorf("orderUsecases.GetActiveOrders() error = %v, expected error %v", err, errExpected)
				}
				return
			}
			if !reflect.DeepEqual(gotOrders, tt.wantOrders) {
				t.Errorf("orderUsecases.GetActiveOrders() = %v, want %v", gotOrders, tt.wantOrders)
			}
		})
	}
}

func TestNewOrderUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exchange := mocks.NewMockExchange(ctrl)
	u := NewOrderUsecase(exchange)
	assert.IsType(t, &orderUsecases{}, u)
	assert.Equal(t, u.(*orderUsecases).exchange, exchange)
	assert.NotNil(t, u.(*orderUsecases).log)
}
