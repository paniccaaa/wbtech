package order

import (
	"context"

	"github.com/paniccaaa/wbtech/internal/model"
)

type Storage interface {
	GetOrder(ctx context.Context, order_uid model.OrderUID) (model.Order, error)
	SaveOrder(ctx context.Context, order model.Order) error
}

type Service struct {
	ordersRepository Storage
}

func NewService(ordersRepo Storage) *Service {
	return &Service{
		ordersRepository: ordersRepo,
	}
}

func (s *Service) GetOrder(ctx context.Context, order_uid model.OrderUID) (model.Order, error) {
	return s.ordersRepository.GetOrder(ctx, order_uid)
}

func (s *Service) SaveOrder(ctx context.Context, order model.Order) error {
	return s.ordersRepository.SaveOrder(ctx, order)
}
