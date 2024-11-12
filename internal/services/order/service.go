package order

import (
	"context"

	"github.com/paniccaaa/wbtech/internal/api/order"
	"github.com/paniccaaa/wbtech/internal/model"
)

type Storage interface {
	GetOrder(ctx context.Context, order_uid model.OrderUID) (model.Order, error)
}

type Service struct {
	ordersRepository Storage
}

func NewService(ordersRepo Storage) order.OrderService {
	return &Service{
		ordersRepository: ordersRepo,
	}
}

func (s *Service) GetOrder(ctx context.Context, order_uid model.OrderUID) (model.Order, error) {
	// TODO:
	return model.Order{}, nil
}
