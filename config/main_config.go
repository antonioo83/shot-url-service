package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseUrl         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	IsUseFileStore  bool
}

func GetConfig() Config {
	const ServerAddress string = ":8080"
	const BaseUrl string = ""
	const FileStoragePath string = "C:\\Users\\jurchenko\\GolandProjects\\short-url-service\\shot-url-service\\data\\database.txt"
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

	cfg.IsUseFileStore = true
	if cfg.FileStoragePath == "" {
		cfg.IsUseFileStore = false
	}

	return cfg
}
