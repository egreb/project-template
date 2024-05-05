package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerPort string `envconfig:"SERVER_PORT"`
	Debug      bool   `envconfig:"DEBUG"`
	Database   DBConfig
}

type DBConfig struct {
	Name     string `envconfig:"DB_NAME"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	Port     string `envconfig:"DB_PORT"`
	Host     string `envconfig:"DB_HOST"`
	SSLMode  string `envconfig:"DB_SSLMODE"`
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed loading .env file: %w", err)
	}
	var c Config
	err = envconfig.Process("", &c)
	if err != nil {
		return nil, fmt.Errorf("failed processing config: %w", err)
	}

	return &c, nil
}
