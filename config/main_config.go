package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseUrl         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	IsUseFileStore  bool
	isInitialized   bool
}

var cfg Config

func GetConfig() Config {
	const ServerAddress string = ":8080"
	const BaseUrl string = ""
	const FileStoragePath string = "C:\\Users\\jurchenko\\GolandProjects\\short-url-service\\shot-url-service\\data\\database.txt"

	if cfg.isInitialized == true {
		return cfg
	}

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "The address of the local server")
	//pflag.Lookup("a").NoOptDefVal = cfg.ServerAddress

	flag.StringVar(&cfg.BaseUrl, "b", cfg.BaseUrl, "Base address of the result short url")
	//pflag.Lookup("b").NoOptDefVal = cfg.BaseUrl

	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "Full filepath to the file storage")
	//pflag.Lookup("f").NoOptDefVal = cfg.FileStoragePath
	flag.Parse()
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

	cfg.isInitialized = true

	return cfg
}
