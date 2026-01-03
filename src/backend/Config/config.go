package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv         string
	Port           string
	DSN            string
	ByPassLogin    string
	API_SECRET_KEY string
}

func LoadConfig() *Config {

	env := os.Getenv("APP_ENV")

	switch env {
	case "development":
		godotenv.Load(".env.development")
	default:
		godotenv.Load(".env.local")
	}

	cfg := &Config{
		AppEnv:      os.Getenv("APP_ENV"),
		Port:        os.Getenv("APP_PORT"),
		DSN:         os.Getenv("DSN"),
		ByPassLogin: os.Getenv("ByPassLogin"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.DSN == "" {
		cfg.DSN = "postgres://postgres:1234@localhost:5432/Project"
	}

	return cfg
}
