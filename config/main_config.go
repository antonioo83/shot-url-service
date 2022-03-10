package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	IsUseFileStore  bool
}

var cfg Config

func GetConfigSettings() Config {
	const ServerAddress string = ":8080"
	const BaseURL string = ""
	const FileStoragePath string = "C:\\Users\\Антон\\GolandProjects\\short-url-service\\shot-url-service\\data\\database.txt"

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "The address of the local server")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base address of the result short url")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "Full filepath to the file storage")
	flag.Parse()
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ServerAddress
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = BaseURL
	}

	cfg.FileStoragePath = FileStoragePath
	cfg.IsUseFileStore = true
	if cfg.FileStoragePath == "" {
		cfg.IsUseFileStore = false
	}

	return cfg
}
