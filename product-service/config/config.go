package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName      string
	AppEnv       string
	AppPort      string
	DatabaseURL  string
	RabbitMQURL  string
	KafkaBrokers string
	KafkaTopic   string
	RedisURL     string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppName:      os.Getenv("APP_NAME"),
		AppEnv:       os.Getenv("APP_ENV"),
		AppPort:      os.Getenv("APP_PORT"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		RabbitMQURL:  os.Getenv("RABBITMQ_URL"),
		KafkaBrokers: os.Getenv("KAFKA_BROKERS"),
		KafkaTopic:   os.Getenv("KAFKA_TOPIC"),
		RedisURL:     os.Getenv("REDIS_URL"),
	}
}
