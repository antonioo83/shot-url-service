package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseUrl       string `env:"BASE_URL"`
}

func GetConfig() Config {
	const ServerAddress string = "http://127.0.0.1:8080"
	const BaseUrl string = "http://127.0.0.1:8080"
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ServerAddress
	}

	if cfg.BaseUrl == "" {
		cfg.BaseUrl = BaseUrl
	}

	return cfg
}
