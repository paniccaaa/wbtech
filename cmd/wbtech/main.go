package main

import (
	"context"
	"log/slog"
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
	log := app.SetupLogger()
	cfg := app.NewConfig()

	db, err := postgres.NewRepository(cfg.DB_URI, log)
	if err != nil {
		log.Error("failed to init db", slog.String("err", err.Error()))
		os.Exit(1)
	}

	schemaClient, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.Kafka.SchemaURI))
	if err != nil {
		log.Error("failed to create schema registry client", slog.String("err", err.Error()))
	}

	orderService, err := order.NewService(db, schemaClient)
	if err != nil {
		log.Error("failed to create order service", slog.String("err", err.Error()))
		os.Exit(1)
	}

	router := app.InitRouter(orderService, log)

	srv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", slog.String("err", err.Error()))
		}
	}()

	log.Info("server started", slog.String("addr", cfg.Server.Addr))

	producer, err := kafka.NewProducer(cfg, schemaClient)
	if err != nil {
		log.Error("failed to create producer", slog.String("err", err.Error()))
	}
	defer producer.Close()

	if err := producer.StartProduce(); err != nil {
		log.Error("failed to produce message", slog.String("err", err.Error()))
		os.Exit(1)
	}

	consumer, err := kafka.NewConsumer(cfg, schemaClient, orderService, log)
	if err != nil {
		log.Error("failed to create consumer", slog.String("err", err.Error()))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Consume(ctx); err != nil {
			log.Error("consumer stopped with error", slog.String("err", err.Error()))
		} else {
			log.Info("Consumer has finished consuming messages")
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)

	<-done

	log.Info("stopping server")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.String("err", err.Error()))
	}

	log.Info("Server stopped.")
}
