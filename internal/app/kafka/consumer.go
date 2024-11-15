package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/paniccaaa/wbtech/internal/app"
	"github.com/paniccaaa/wbtech/internal/model"
)

type Consumer struct {
	client       *kafka.Consumer
	cfgKafka     app.Kafka
	schemaClient schemaregistry.Client
}

func NewConsumer(cfg app.Config, schemaClient schemaregistry.Client) (*Consumer, error) {
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
	}, nil
}

func (c *Consumer) ListenAndConsume() error {
	partitions, err := c.client.GetMetadata(&c.cfgKafka.Topic, false, 10000)
	if err != nil {
		return fmt.Errorf("get metadata: %w", err)
	}

	for _, partition := range partitions.Topics[c.cfgKafka.Topic].Partitions {
		go func() {
			if err := c.consumePartitionMessages(partition.ID); err != nil {
				log.Printf("consume partition messages: %v", err)
			}
		}()

	}

	return nil
}

func (c *Consumer) consumePartitionMessages(partition int32) error {
	err := c.client.Assign([]kafka.TopicPartition{
		{
			Topic:     &c.cfgKafka.Topic,
			Partition: partition,
		},
	})
	if err != nil {
		return fmt.Errorf("assign partition: %w", err)
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

		switch e := ev.(type) {
		case *kafka.Message:
			var value model.Order
			err := deser.DeserializeInto(*e.TopicPartition.Topic, e.Value, &value)
			if err != nil {
				log.Printf("Failed to deserialize message: %s", err)
			} else {
				log.Printf("Received message from partition %d: %+v\n", partition, value)
			}
		case kafka.Error:
			log.Printf("Error: %v\n", e)
		default:
			log.Printf("Ignored %v\n", e)
		}
	}
}
