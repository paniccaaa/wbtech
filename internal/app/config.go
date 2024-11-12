package app

import "os"

type Config struct {
	DB_URI string `env:"DB_URI"`
	Kafka  Kafka  `env:"KAFKA"`
}

type Kafka struct {
	URI   string `env:"URI"`
	Topic string `env:"TOPIC"`
}

func NewConfig() Config {
	return Config{
		DB_URI: os.Getenv("DB_URI"),
		Kafka: Kafka{
			URI:   "localhost:9092",
			Topic: "orders-topic",
		},
	}
}
