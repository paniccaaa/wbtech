package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/paniccaaa/wbtech/internal/api/kafka"
	"github.com/paniccaaa/wbtech/internal/app"
	"github.com/paniccaaa/wbtech/internal/repository/postgres"
	"github.com/paniccaaa/wbtech/internal/services/order"
)

func main() {
	cfg := app.NewConfig()

	db, err := postgres.NewRepository(cfg.DB_URI)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	orderService := order.NewService(db)
	router := app.InitRouter(orderService)

	srv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("failed to start server: %v", err)
		}
	}()

	log.Println("Server started")

	schemaClient, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.Schema_URI))
	if err != nil {
		log.Fatalf("failed to create schema registry client: %v", err)
	}

	producer, err := kafka.NewProducer(cfg, schemaClient)
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}
	defer producer.Close()

	consumer, err := kafka.NewConsumer(cfg, schemaClient, orderService)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	if err := producer.StartProduce(); err != nil {
		log.Fatalf("failed to produce messages: %v", err)
	}

	log.Println("Producer has finished sending messages to Kafka.")

	go func() {
		if err := consumer.ListenAndConsume(); err != nil {
			log.Printf("Consumer error: %v", err)
		} else {
			log.Println("Consumer has finished consuming messages.")
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done

	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to stop server: %v", err)
	}

	log.Println("Server gracefully stopped.")
}
