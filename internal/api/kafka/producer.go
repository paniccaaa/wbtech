package kafka

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/paniccaaa/wbtech/internal/app"
	"github.com/paniccaaa/wbtech/internal/model"
)

type Producer struct {
	client       *kafka.Producer
	schemaClient schemaregistry.Client
	cfgKafka     app.Kafka
}

func NewProducer(cfg *app.Config, schemaClient schemaregistry.Client) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.URI,
	})
	if err != nil {
		return nil, fmt.Errorf("new producer: %w", err)
	}

	return &Producer{
		client:       p,
		schemaClient: schemaClient,
		cfgKafka:     cfg.Kafka,
	}, nil
}

func (p *Producer) Close() {
	p.client.Close()
}

func (p *Producer) StartProduce() error {
	ser, err := jsonschema.NewSerializer(p.schemaClient, serde.ValueSerde, jsonschema.NewSerializerConfig())
	if err != nil {
		return fmt.Errorf("create serializer: %w", err)
	}

	fileData, err := os.ReadFile("orders.json")
	if err != nil {
		return fmt.Errorf("read JSON file: %w", err)
	}

	var orders []model.Order
	if err := json.Unmarshal(fileData, &orders); err != nil {
		return fmt.Errorf("unmarshal JSON data: %w", err)
	}

	for _, order := range orders {
		if err := p.ProduceMessage(ser, order); err != nil {
			return fmt.Errorf("produce message: %w", err)
		}
	}

	p.client.Flush(3 * 1000)
	return nil
}

func (p *Producer) ProduceMessage(ser *jsonschema.Serializer, order model.Order) error {
	data, err := ser.Serialize(p.cfgKafka.Topic, &order)
	if err != nil {
		return fmt.Errorf("marshal order to JSON: %w", err)
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.cfgKafka.Topic, Partition: kafka.PartitionAny},
		Key:            []byte(order.OrderUID),
		Value:          data,
	}

	err = p.client.Produce(message, nil)
	if err != nil {
		return fmt.Errorf("produce message: %w", err)
	}

	return nil
}
