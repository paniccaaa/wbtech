package order

import (
	"context"
	"reflect"
	"testing"

	"github.com/paniccaaa/wbtech/internal/model"
	"github.com/paniccaaa/wbtech/internal/services/order/mocks"
)

func TestService_GetOrder(t *testing.T) {
	type args struct {
		ctx      context.Context
		orderUID model.OrderUID
	}
	tests := []struct {
		name       string
		args       args
		mockReturn struct {
			order model.Order
			err   error
		}
		want    model.Order
		wantErr bool
	}{
		{
			name: "successful get",
			args: args{
				ctx:      context.Background(),
				orderUID: "b563feb7b2b84b6test",
			},
			mockReturn: struct {
				order model.Order
				err   error
			}{
				order: model.Order{
					OrderUID: "b563feb7b2b84b6test",
				},
				err: nil,
			},
			want: model.Order{
				OrderUID: "b563feb7b2b84b6test",
			},
			wantErr: false,
		},
		{
			name: "order not found",
			args: args{
				ctx:      context.Background(),
				orderUID: "non-existent order uid",
			},
			mockReturn: struct {
				order model.Order
				err   error
			}{
				order: model.Order{},
				err:   model.ErrOrderNotFound,
			},
			want:    model.Order{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)

			storage.
				On("GetOrder", tt.args.ctx, tt.args.orderUID).
				Return(tt.mockReturn.order, tt.mockReturn.err).
				Once()

			s := &Service{
				ordersRepository: storage,
			}

			got, err := s.GetOrder(tt.args.ctx, tt.args.orderUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ProcessKafkaMessage(t *testing.T) {
	type args struct {
		ctx     context.Context
		topic   string
		message []byte
	}
	type mockReturn struct {
		order model.Order
		err   error
	}
	tests := []struct {
		name       string
		args       args
		mockReturn mockReturn
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "successful process",
			args: args{
				ctx:     context.Background(),
				topic:   "orders",
				message: []byte(`{"OrderUID": "b563feb7b2b84b6test"}`),
			},
			mockReturn: mockReturn{
				err: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)
			deser := mocks.NewDeserializer(t)

			deser.
				On("DeserializeInto", tt.args.topic, tt.args.message, &model.Order{}).
				Return(tt.mockReturn.err).
				Once()

			storage.
				On("SaveOrder", tt.args.ctx, tt.mockReturn.order).
				Return(tt.mockReturn.err).
				Once()

			s := &Service{
				ordersRepository: storage,
				deser:            deser,
			}

			if err := s.ProcessKafkaMessage(tt.args.ctx, tt.args.topic, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Service.ProcessKafkaMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
