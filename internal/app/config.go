package app

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB_URI string `yaml:"db_uri"`
	Kafka  Kafka  `yaml:"kafka"`
	Server Server `yaml:"srv"`
}

type Kafka struct {
	URI       string `yaml:"uri"`
	Topic     string `yaml:"topic"`
	SchemaURI string `yaml:"schema_uri"`
}

type Server struct {
	Addr string `yaml:"addr"`
}

func NewConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("config path is required")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return &cfg
}
