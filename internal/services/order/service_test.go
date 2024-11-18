package order

import (
	"context"
	"testing"

	"github.com/paniccaaa/wbtech/internal/model"
	mock_order "github.com/paniccaaa/wbtech/internal/services/order/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type serviceTestDeps struct {
	Service *Service
	ctx     context.Context
	storage *mock_order.MockStorage
}

func getTestDeps(t *testing.T) *serviceTestDeps {
	ctrl := gomock.NewController(t)
	storage := mock_order.NewMockStorage(ctrl)

	return &serviceTestDeps{
		Service: &Service{
			ordersRepository: storage,
		},
		ctx:     context.Background(),
		storage: storage,
	}
}

func TestService_GetOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUID model.OrderUID
	}

	type result struct {
		order model.Order
		err   error
	}

	type testCase struct {
		name   string
		args   args
		mocks  func(tc testCase, deps *serviceTestDeps)
		result result
	}

	testCases := []testCase{
		{
			name: "success",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.storage.EXPECT().
					GetOrder(gomock.Any(), tc.args.orderUID).
					Return(model.Order{OrderUID: tc.args.orderUID}, nil)
			},
			args: args{
				orderUID: model.OrderUID("b563feb7b2b84b6test13"),
			},
			result: result{
				order: model.Order{OrderUID: "b563feb7b2b84b6test13"},
				err:   nil,
			},
		},
		{
			name: "order not found",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.storage.EXPECT().
					GetOrder(gomock.Any(), tc.args.orderUID).
					Return(model.Order{}, model.ErrOrderNotFound)
			},
			args: args{
				orderUID: model.OrderUID("456"),
			},
			result: result{
				order: model.Order{},
				err:   model.ErrOrderNotFound,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := getTestDeps(t)

			tc.mocks(tc, deps)

			order, err := deps.Service.GetOrder(deps.ctx, tc.args.orderUID)
			require.ErrorIs(t, err, tc.result.err)
			require.Equal(t, tc.result.order, order)
		})
	}
}

func TestService_SaveOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		order model.Order
	}

	type result struct {
		err error
	}

	type testCase struct {
		name   string
		args   args
		mocks  func(tc testCase, deps *serviceTestDeps)
		result result
	}

	testCases := []testCase{
		{
			name: "success",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.storage.EXPECT().
					SaveOrder(gomock.Any(), tc.args.order).
					Return(nil)
			},
			args: args{
				order: model.Order{
					OrderUID:    "123",
					TrackNumber: "TN123",
					CustomerID:  "customer1",
					Delivery: model.Delivery{
						Name:    "John Doe",
						Phone:   "1234567890",
						Zip:     "12345",
						City:    "City",
						Address: "123 St.",
						Region:  "Region",
						Email:   "johndoe@example.com",
					},
					Payment: model.Payment{
						Transaction:  "txn123",
						RequestID:    "req123",
						Currency:     "USD",
						Provider:     "Provider1",
						Amount:       1000,
						PaymentDT:    1616161616,
						Bank:         "Bank1",
						DeliveryCost: 50,
						GoodsTotal:   950,
						CustomFee:    10,
					},
					Items: []model.Item{
						{
							ChrtID:      1,
							TrackNumber: "TN123",
							Price:       500,
							Rid:         "RID123",
							Name:        "Item1",
							Sale:        10,
							Size:        "L",
							TotalPrice:  450,
							NmID:        1234,
							Brand:       "Brand1",
							Status:      1,
						},
					},
				},
			},
			result: result{
				err: nil,
			},
		},
		// {
		// 	name: "save order failed",
		// 	mocks: func(tc testCase, deps *serviceTestDeps) {
		// 		deps.storage.EXPECT().
		// 			SaveOrder(gomock.Any(), tc.args.order)
		// 		// Return(model.ErrDatabase)
		// 	},
		// 	args: args{
		// 		order: model.Order{OrderUID: "123"},
		// 	},
		// 	result: result{
		// 		// err: model.ErrDatabase,
		// 	},
		// },
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := getTestDeps(t)

			tc.mocks(tc, deps)

			err := deps.Service.SaveOrder(deps.ctx, tc.args.order)
			require.ErrorIs(t, err, tc.result.err)
		})
	}
}
