package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName     string
	AppEnv      string
	AppPort     string
	DatabaseURL string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppName:     os.Getenv("APP_NAME"),
		AppEnv:      os.Getenv("APP_ENV"),
		AppPort:     os.Getenv("APP_PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
