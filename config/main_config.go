package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress       string `env:"SERVER_ADDRESS"`
	BaseURL             string `env:"BASE_URL"`
	FileStoragePath     string `env:"FILE_STORAGE_PATH"`
	UserFileStoragePath string `env:"USER_FILE_STORAGE_PATH"`
	IsUseFileStore      bool
	DatabaseDsn         string `env:"DATABASE_DSN"`
}

var cfg Config

func GetConfigSettings() Config {
	const ServerAddress string = ":8080"
	const BaseURL string = ""
	//const FileStoragePath string = "..\\data\\short_url_database.txt"
	const UserFileStoragePath string = "user_database.txt"
	const DatabaseDSN = "postgres://postgres:433370@localhost:5433/postgres"

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "The address of the local server")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base address of the result short url")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "Full filepath to the file storage")
	flag.StringVar(&cfg.UserFileStoragePath, "fu", cfg.UserFileStoragePath, "Full filepath to the user file storage")
	flag.StringVar(&cfg.DatabaseDsn, "d", cfg.DatabaseDsn, "Database port")
	flag.Parse()
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ServerAddress
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = BaseURL
	}

	cfg.IsUseFileStore = true
	if cfg.FileStoragePath == "" {
		cfg.IsUseFileStore = false
	} else {
		cfg.UserFileStoragePath = UserFileStoragePath
	}

	if cfg.DatabaseDsn == "" {
		cfg.DatabaseDsn = DatabaseDSN
	}

	return cfg
}
