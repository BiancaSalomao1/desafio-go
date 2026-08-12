package config

/*
struct Config

- armazenar configurações da aplicação.

Métodos:
- Load()
*/

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName      string
	AppEnv       string
	AppPort      string
	ServerPort   string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn string
}

func Load() *Config {

	// Carrega o arquivo .env caso exista.
	if err := godotenv.Load(); err != nil {
		log.Println(".env não encontrado. Utilizando variáveis de ambiente.")
	}

	return &Config{
		AppName:      os.Getenv("APP_NAME"),
		AppEnv:       os.Getenv("APP_ENV"),
		AppPort:      os.Getenv("APP_PORT"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		JWTExpiresIn: os.Getenv("JWT_EXPIRES_IN"),
	}
}
