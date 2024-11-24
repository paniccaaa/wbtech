package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	api "github.com/paniccaaa/wbtech/internal/api/kafka"
	"github.com/paniccaaa/wbtech/internal/app"
	"github.com/paniccaaa/wbtech/internal/services/order"
)

func main() {
	log := app.SetupLogger()
	cfg := app.NewConfig()

	db := app.SetupDB(cfg.DB_URI, log)

	deser, schemaClient := app.SetupDeserializer(cfg, log)
	orderService := order.NewService(db, deser, log)

	srv := app.SetupServer(cfg, orderService, log)
	app.StartServer(srv, log)

	log.Info("server started", slog.String("addr", cfg.Server.Addr))

	producer, err := api.NewProducer(cfg, schemaClient)
	if err != nil {
		log.Error("failed to create producer", slog.String("err", err.Error()))
	}
	defer producer.Close()

	if err := producer.StartProduce(); err != nil {
		log.Error("failed to produce message", slog.String("err", err.Error()))
		os.Exit(1)
	}

	c := app.SetupConsumer(cfg, log)
	consumer := api.NewConsumer(cfg, c, orderService, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Consume(ctx); err != nil {
			log.Error("consumer stopped with error", slog.String("err", err.Error()))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)

	<-done

	app.ShutdownServer(srv, log)
}
