package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/paniccaaa/wbtech/internal/repository/postgres"
	"github.com/paniccaaa/wbtech/internal/services/order"
)

func SetupDB(dbURI string, log *slog.Logger) order.Storage {
	db, err := postgres.NewRepository(dbURI, log)
	if err != nil {
		log.Error("failed to init db", slog.String("err", err.Error()))
		os.Exit(1)
	}

	return db
}

func SetupServer(cfg *Config, orderService *order.Service, log *slog.Logger) *http.Server {
	router := InitRouter(orderService, log)

	return &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: router,
	}
}

func StartServer(srv *http.Server, log *slog.Logger) {
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", slog.String("err", err.Error()))
		}
	}()
}

func ShutdownServer(srv *http.Server, log *slog.Logger) {
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.String("err", err.Error()))
	}

	log.Info("Server stopped.")
}

func SetupConsumer(cfg *Config, log *slog.Logger) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.URI,
		"group.id":          "group-id",
		"auto.offset.reset": "earliest",
		"security.protocol": "PLAINTEXT",
	})
	if err != nil {
		log.Error("create consumer", slog.String("err", err.Error()))
		os.Exit(1)
	}

	if err := c.Subscribe(cfg.Kafka.Topic, nil); err != nil {
		log.Error("subscribe to topic", slog.String("err", err.Error()))
		os.Exit(1)
	}

	return c
}

func SetupDeserializer(cfg *Config, log *slog.Logger) (*jsonschema.Deserializer, schemaregistry.Client) {
	schemaClient, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.Kafka.SchemaURI))
	if err != nil {
		log.Error("failed to create schema registry client", slog.String("err", err.Error()))
		os.Exit(1)
	}

	deser, err := jsonschema.NewDeserializer(
		schemaClient,
		serde.ValueSerde,
		jsonschema.NewDeserializerConfig(),
	)
	if err != nil {
		log.Error("failed to create deserializer", slog.String("err", err.Error()))
		os.Exit(1)
	}

	return deser, schemaClient
}
