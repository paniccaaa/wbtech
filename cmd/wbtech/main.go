package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/paniccaaa/wbtech/internal/app"
	"github.com/paniccaaa/wbtech/internal/model"
	"github.com/paniccaaa/wbtech/internal/repository/postgres"
	"github.com/paniccaaa/wbtech/internal/services/order"
)

// func producer(cfg app.Config, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Kafka.URI})
// 	if err != nil {
// 		log.Fatalf("failed to create producer: %v", err)
// 	}
// 	defer p.Close()

// 	fileData, err := os.ReadFile("orders.json")
// 	if err != nil {
// 		log.Fatalf("failed to read JSON file: %v", err)
// 	}

// 	var orders []model.Order
// 	if err := json.Unmarshal(fileData, &orders); err != nil {
// 		log.Fatalf("failed to unmarshall JSON data: %v", err)
// 	}

// 	message := &kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &cfg.Kafka.Topic, Partition: kafka.PartitionAny},
// 		Key:            []byte("order-key"),
// 		Value:          []byte(`123`),
// 	}

// 	err = p.Produce(message, nil)
// 	if err != nil {
// 		log.Fatalf("failed to produce message: %v", err)
// 	}

// 	p.Flush(15 * 1000)
// 	fmt.Println("Message sent to Kafka topic")
// }

func producer(cfg app.Config, wg *sync.WaitGroup) {
	defer wg.Done()

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Kafka.URI})
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}
	defer p.Close()

	fileData, err := os.ReadFile("orders.json")
	if err != nil {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	var orders []model.Order
	err = json.Unmarshal(fileData, &orders)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	for _, order := range orders {
		_, err := json.Marshal(order)
		if err != nil {
			log.Fatalf("failed to marshal order to JSON: %v", err)
		}

		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &cfg.Kafka.Topic, Partition: kafka.PartitionAny},
			Key:            []byte(order.OrderUID),
			Value:          []byte(`123`),
		}

		err = p.Produce(message, nil)
		if err != nil {
			log.Fatalf("failed to produce message: %v", err)
		}
	}

	p.Flush(15 * 1000)
	fmt.Println("All messages sent to Kafka topic")
}

func consumer(cfg app.Config, wg *sync.WaitGroup) {
	defer wg.Done()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.URI,
		"group.id":          "group-id",
		"auto.offset.reset": "latest",
		"security.protocol": "PLAINTEXT",
	})
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	defer func(c *kafka.Consumer) {
		if err := c.Close(); err != nil {
			log.Fatalf("failed to close consumer instance: %v", err)
		}
	}(c)

	err = c.SubscribeTopics([]string{cfg.Kafka.Topic}, nil)
	if err != nil {
		log.Fatalf("failed to subscribe to topics: %v", err)
	}

	log.Println("Starting Kafka consumer...")

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			log.Fatalf("could not read message: %v", err)
		}
		fmt.Printf("Received message at offset %d: %s = %s\n", msg.TopicPartition.Offset, string(msg.Key), string(msg.Value))
	}
}

func main() {
	cfg := app.NewConfig()

	var wg sync.WaitGroup

	wg.Add(2)
	go producer(cfg, &wg)
	go consumer(cfg, &wg)
	wg.Wait()

	db := postgres.NewRepository(cfg.DB_URI)

	orderService := order.NewService(db)

	router := app.InitRouter(orderService)

	srv := &http.Server{
		Addr:    "0.0.0.0:8089",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	log.Println("start server")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to stop server: %v", err)
		return
	}
}
