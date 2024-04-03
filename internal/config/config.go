package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server   `yaml:"server"`
		RabbitMQ `yaml:"rabbitmq"`
		PG
		Log `yaml:"logger"`
	}
	Server struct {
		Port         string `yaml:"port"`
		ReadTimeout  int    `yaml:"read_timeout"`
		WriteTimeout int    `yaml:"write_timeout"`
	}
	RabbitMQ struct {
		URL          string
		Exchange     string `yaml:"exchange"`
		ExchangeType string `yaml:"exchange_type"`
		Queue        string `yaml:"queue"`
	}

	PG struct {
		URL string
	}
	Log struct {
		Level string `yaml:"log_level"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	cfg.PG.URL = os.Getenv("PG_URL")
	cfg.RabbitMQ.URL = os.Getenv("RABBITMQ_URL")
	fmt.Println(cfg.PG.URL)
	return cfg, nil
}
