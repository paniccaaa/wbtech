package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/paniccaaa/wbtech/internal/app"
	"github.com/paniccaaa/wbtech/internal/model"
)

type SaveProvider interface {
	SaveOrder(ctx context.Context, order model.Order) error
}

type Consumer struct {
	client       *kafka.Consumer
	cfgKafka     app.Kafka
	schemaClient schemaregistry.Client
	order        SaveProvider
}

func NewConsumer(cfg app.Config, schemaClient schemaregistry.Client, order SaveProvider) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.URI,
		"group.id":          "group-id",
		"auto.offset.reset": "earliest",
		"security.protocol": "PLAINTEXT",
	})
	if err != nil {
		return nil, fmt.Errorf("create consumer: %w", err)
	}

	return &Consumer{
		client:       c,
		schemaClient: schemaClient,
		cfgKafka:     cfg.Kafka,
		order:        order,
	}, nil
}

func (c *Consumer) ListenAndConsume() error {
	partitions, err := c.client.GetMetadata(&c.cfgKafka.Topic, false, 10000)
	if err != nil {
		return fmt.Errorf("get metadata: %w", err)
	}

	for _, partition := range partitions.Topics[c.cfgKafka.Topic].Partitions {
		go func(partitionID int32) {
			if err := c.consumePartitionMessages(partitionID); err != nil {
				log.Printf("consume partition messages: %v", err)
			}
		}(partition.ID)
	}

	return nil
}

func (c *Consumer) consumePartitionMessages(partition int32) error {
	err := c.client.Subscribe(c.cfgKafka.Topic, nil)
	if err != nil {
		return fmt.Errorf("subscribe to topic: %w", err)
	}

	deser, err := jsonschema.NewDeserializer(
		c.schemaClient,
		serde.ValueSerde,
		jsonschema.NewDeserializerConfig(),
	)
	if err != nil {
		return fmt.Errorf("create deserializer: %w", err)
	}

	for {
		ev := c.client.Poll(100)
		if ev == nil {
			continue
		}

		if err := c.processEvent(ev, deser, partition); err != nil {
			log.Printf("Failed to process event: %v", err)
		}
	}
}

func (c *Consumer) processEvent(ev kafka.Event, deser *jsonschema.Deserializer, partition int32) error {
	switch e := ev.(type) {
	case *kafka.Message:
		return c.handleMessage(e, deser, partition)
	case kafka.Error:
		log.Printf("Kafka error: %v", e)
		return nil
	default:
		log.Printf("Ignored unexpected event: %v", e)
		return nil
	}
}

func (c *Consumer) handleMessage(msg *kafka.Message, deser *jsonschema.Deserializer, partition int32) error {
	var order model.Order
	if err := deser.DeserializeInto(*msg.TopicPartition.Topic, msg.Value, &order); err != nil {
		return fmt.Errorf("failed to deserialize message: %w", err)
	}

	if err := c.order.SaveOrder(context.Background(), order); err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	log.Printf("Received and saved message from partition %d with order_uid=%v", partition, order.OrderUID)

	return nil
}
