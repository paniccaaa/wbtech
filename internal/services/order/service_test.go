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

func TestService_SaveOrder(t *testing.T) {
	type args struct {
		ctx   context.Context
		order model.Order
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "successful save",
			args: args{
				ctx: context.Background(),
				order: model.Order{
					OrderUID: "b563feb7b2b84b6test",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)

			storage.
				On("SaveOrder", tt.args.ctx, tt.args.order).
				Return().
				Once()

			s := &Service{
				ordersRepository: storage,
			}

			if err := s.SaveOrder(tt.args.ctx, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("Service.SaveOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
