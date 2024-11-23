package order

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/paniccaaa/wbtech/internal/model"
)

//go:generate mockery --name Storage
type Storage interface {
	GetOrder(ctx context.Context, orderUID model.OrderUID) (model.Order, error)
	SaveOrder(ctx context.Context, order model.Order) error
}

type Service struct {
	ordersRepository Storage
	deser            *jsonschema.Deserializer
	schemaClient     schemaregistry.Client
}

func NewService(ordersRepo Storage, schemaClient schemaregistry.Client) (*Service, error) {
	deser, err := jsonschema.NewDeserializer(
		schemaClient,
		serde.ValueSerde,
		jsonschema.NewDeserializerConfig(),
	)
	if err != nil {
		return nil, fmt.Errorf("create deserializer: %w", err)
	}

	return &Service{
		ordersRepository: ordersRepo,
		deser:            deser,
		schemaClient:     schemaClient,
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

	log.Printf("Successfully processed and saved order: %v", order.OrderUID)

	return nil
}
