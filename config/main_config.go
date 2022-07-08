// Package config This module is intended for service configuration.
package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

// Config Configuration settings.
type Config struct {
	ServerAddress       string         `env:"SERVER_ADDRESS" json:"server_address,omitempty"`                 // The address of the local server.
	BaseURL             string         `env:"BASE_URL" json:"base_url,omitempty"`                             // Base address of the result short url
	FileStoragePath     string         `env:"FILE_STORAGE_PATH" json:"file_storage_path,omitempty"`           // Filepath to the shot url file storage.
	UserFileStoragePath string         `env:"USER_FILE_STORAGE_PATH" json:"user_file_storage_path,omitempty"` // Filepath to the user file storage.
	IsUseFileStore      bool           // Is use file store ?
	DatabaseDsn         string         `env:"DATABASE_DSN" json:"database_dsn,omitempty"` // Database connection string.
	IsUseDatabase       bool           // Is use database ?
	FilepathToDBDump    string         // Filepath to the SQL dump for initialization database.
	EnableHTTPS         bool           `env:"ENABLE_HTTPS" json:"enable_https,omitempty"` // Enable HTTPS connection.
	ConfigFilePath      string         `env:"CONFIG" json:"config_file_path,omitempty"`   // Filename of the server configurations.
	Auth                Auth           // User Authorization settings
	DeleteShotURL       DeleteShortURL // Settings of deleting short url rows from database
}

// Auth User Authorization settings.
type Auth struct {
	Alg       string        // Type of encryption algorithm.
	CookieTTL time.Duration // Time of storage cookie.
	SignKey   []byte        // Authorization secret token name.
	TokenName string        // Authorization token name.
}

// DeleteShortURL Settings of deleting short url rows from database.
type DeleteShortURL struct {
	WorkersCount int // Workers count.
	ChunkLength  int // Length of rows chunk.
}

var cfg Config

// GetConfigSettings returns configuration settings.
func GetConfigSettings(configFromFile *Config) Config {
	const (
		serverAddress       string = ":8080"
		baseURL             string = ""
		userFileStoragePath string = "user_database.txt"
		authEncodeAlgorithm        = "HS256"
		authRememberMeTime         = 60 * 30 * time.Second
		authSignKey                = "secret"
		authTokenName              = "token"
	)

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "The address of the local server")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base address of the result short url")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "Filepath to the file storage")
	flag.StringVar(&cfg.UserFileStoragePath, "fu", cfg.UserFileStoragePath, "Filepath to the user file storage")
	flag.StringVar(&cfg.DatabaseDsn, "d", cfg.DatabaseDsn, "Database port")
	flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "Enable HTTPS connection")
	flag.StringVar(&cfg.ConfigFilePath, "c", cfg.ConfigFilePath, "Filename of the server configurations")
	flag.Parse()

	if configFromFile != nil {
		if cfg.ServerAddress == "" {
			cfg.ServerAddress = configFromFile.ServerAddress
		}
		if cfg.BaseURL == "" {
			cfg.BaseURL = configFromFile.BaseURL
		}
		if cfg.FileStoragePath == "" {
			cfg.FileStoragePath = configFromFile.FileStoragePath
		}
		if cfg.UserFileStoragePath == "" {
			cfg.UserFileStoragePath = configFromFile.UserFileStoragePath
		}
		if cfg.DatabaseDsn == "" {
			cfg.DatabaseDsn = configFromFile.DatabaseDsn
		}
		if !cfg.EnableHTTPS {
			cfg.EnableHTTPS = configFromFile.EnableHTTPS
		}
	}

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = serverAddress
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = baseURL
	}

	cfg.IsUseFileStore = true
	if cfg.FileStoragePath == "" {
		cfg.IsUseFileStore = false
	} else {
		cfg.UserFileStoragePath = userFileStoragePath
	}

	cfg.IsUseDatabase = true
	if cfg.DatabaseDsn == "" {
		cfg.IsUseDatabase = false
	}

	cfg.Auth.Alg = authEncodeAlgorithm
	cfg.Auth.CookieTTL = authRememberMeTime
	cfg.Auth.SignKey = []byte(authSignKey)
	cfg.Auth.TokenName = authTokenName

	cfg.DeleteShotURL.ChunkLength = 10
	cfg.DeleteShotURL.WorkersCount = 1

	return cfg
}
