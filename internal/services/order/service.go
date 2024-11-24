package order

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/paniccaaa/wbtech/internal/model"
)

//go:generate mockery --name Storage
type Storage interface {
	GetOrder(ctx context.Context, orderUID model.OrderUID) (model.Order, error)
	SaveOrder(ctx context.Context, order model.Order) error
}

//go:generate mockery --name Deserializer
type Deserializer interface {
	DeserializeInto(topic string, message []byte, v interface{}) error
}

type Service struct {
	ordersRepository Storage
	deser            Deserializer
	log              *slog.Logger
}

func NewService(ordersRepo Storage, deser Deserializer, log *slog.Logger) (*Service, error) {
	return &Service{
		ordersRepository: ordersRepo,
		deser:            deser,
		log:              log,
	}, nil
}

func (s *Service) GetOrder(ctx context.Context, orderUID model.OrderUID) (model.Order, error) {
	return s.ordersRepository.GetOrder(ctx, orderUID)
}

func (s *Service) ProcessKafkaMessage(ctx context.Context, topic string, message []byte) error {
	var order model.Order
	if err := s.deser.DeserializeInto(topic, message, &order); err != nil {
		return fmt.Errorf("deserialize message: %w", err)
	}

	if err := s.ordersRepository.SaveOrder(ctx, order); err != nil {
		return fmt.Errorf("save order: %w", err)
	}

	s.log.Info("Successfully processed and saved order", slog.String("order_uid", string(order.OrderUID)))

	return nil
}
