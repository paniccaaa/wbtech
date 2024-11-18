package app

type Config struct {
	DB_URI     string `env:"DB_URI"`
	Kafka      Kafka  `env:"KAFKA"`
	Server     Server `env:"SERVER"`
	Schema_URI string `env:"SCHEMA_URI"`
}

type Kafka struct {
	URI   string `env:"URI"`
	Topic string `env:"TOPIC"`
}

type Server struct {
	Addr string `env:"ADDR"`
}

func NewConfig() Config {
	return Config{
		DB_URI: "postgres://wbuser:wbpassword@localhost:5435/postgres?sslmode=disable",
		Kafka: Kafka{
			URI:   "localhost:9092",
			Topic: "orders",
		},
		Server: Server{
			Addr: "localhost:8089",
		},
		Schema_URI: "http://localhost:8081",
	}
}
