package kafka

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/paniccaaa/wbtech/internal/app"
)

type MessageHandler interface {
	ProcessKafkaMessage(ctx context.Context, topic string, message []byte) error
}

type Consumer struct {
	client       *kafka.Consumer
	cfgKafka     app.Kafka
	schemaClient schemaregistry.Client
	handler      MessageHandler
	log          *slog.Logger
}

func NewConsumer(cfg *app.Config, schemaClient schemaregistry.Client, handler MessageHandler, log *slog.Logger) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.URI,
		"group.id":          "group-id",
		"auto.offset.reset": "earliest",
		"security.protocol": "PLAINTEXT",
	})
	if err != nil {
		return nil, fmt.Errorf("create consumer: %w", err)
	}

	if err := c.Subscribe(cfg.Kafka.Topic, nil); err != nil {
		return nil, fmt.Errorf("subscribe to topic: %w", err)
	}

	return &Consumer{
		client:       c,
		schemaClient: schemaClient,
		cfgKafka:     cfg.Kafka,
		handler:      handler,
		log:          log,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context) error {
	for {
		ev := c.client.Poll(100)
		if ev == nil {
			continue
		}

		if err := c.processEvent(ctx, ev); err != nil {
			c.log.Error("process event", slog.String("err", err.Error()))
		}
	}
}
func (c *Consumer) processEvent(ctx context.Context, ev kafka.Event) error {
	switch e := ev.(type) {
	case *kafka.Message:
		return c.handler.ProcessKafkaMessage(ctx, *e.TopicPartition.Topic, e.Value)
	case kafka.Error:
		c.log.Error("Kafka error", slog.String("err", e.Error()))
		return nil
	default:
		c.log.Info("Ignored unexpected event", slog.String("err", e.String()))
		return nil
	}
}
