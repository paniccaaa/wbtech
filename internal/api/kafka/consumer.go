package kafka

import (
	"context"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/paniccaaa/wbtech/internal/app"
)

//go:generate mockery --name MessageHandler
type MessageHandler interface {
	ProcessKafkaMessage(ctx context.Context, topic string, message []byte) error
}

//go:generate mockery --name Poller
type Poller interface {
	Poll(timeoutMs int) (event kafka.Event)
}

type Consumer struct {
	client   Poller
	cfgKafka app.Kafka
	handler  MessageHandler
	log      *slog.Logger
}

func NewConsumer(cfg *app.Config, client Poller, handler MessageHandler, log *slog.Logger) *Consumer {
	return &Consumer{
		client:   client,
		cfgKafka: cfg.Kafka,
		handler:  handler,
		log:      log,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			ev := c.client.Poll(100)
			if ev == nil {
				continue
			}

			if err := c.processEvent(ctx, ev); err != nil {
				c.log.Error("process event", slog.String("err", err.Error()))
			}
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
