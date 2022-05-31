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
	ServerAddress       string        `env:"SERVER_ADDRESS"`         // The address of the local server.
	BaseURL             string        `env:"BASE_URL"`               // Base address of the result short url
	FileStoragePath     string        `env:"FILE_STORAGE_PATH"`      // Full filepath to the shot url file storage.
	UserFileStoragePath string        `env:"USER_FILE_STORAGE_PATH"` // Full filepath to the user file storage.
	IsUseFileStore      bool          // Is use file store ?
	DatabaseDsn         string        `env:"DATABASE_DSN"` // Database connection string.
	IsUseDatabase       bool          // Is use database ?
	FilepathToDBDump    string        // Filepath to the SQL dump for initialization database.
	Auth                Auth          //
	DeleteShotURL       DeleteShotURL //
}

// Auth User Authorization settings.
type Auth struct {
	Alg            string        //Type of encryption algorithm.
	RememberMeTime time.Duration //Time of storage cookie.
	SignKey        []byte        // Authorization secret token name.
	TokenName      string        // Authorization token name.
}

//DeleteShotURL Settings of deleting short url rows from database.
type DeleteShotURL struct {
	WorkersCount int //Workers count.
	ChunkLength  int //Length of rows chunk.
}

var cfg Config

// GetConfigSettings — returns configuration settings.
func GetConfigSettings() Config {
	const ServerAddress string = ":8080" // The address of the local server.
	const BaseURL string = ""            //Base address of the result short url
	//const FileStoragePath string = "..\\data\\short_url_database.txt"
	const UserFileStoragePath string = "user_database.txt" // Full filepath to the user file storage.
	//const DatabaseDSN = "postgres://postgres:433370@localhost:5433/postgres"
	const AuthEncodeAlgorithm = "HS256"
	const AuthRememberMeTime = 60 * 30 * time.Second //Cookie storage time.
	const AuthSignKey = "secret"                     // Authorization secret token name.
	const AuthTokenName = "token"                    // Authorization token name.

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

	cfg.IsUseDatabase = true
	if cfg.DatabaseDsn == "" {
		cfg.IsUseDatabase = false
	} //else {
	// GitHub test environment doesn't understand sql dump as file.
	//cfg.FilepathToDBDump, _ = os.Getwd()
	//cfg.FilepathToDBDump += "\\migrations\\create_tables.sql"
	//}

	cfg.Auth.Alg = AuthEncodeAlgorithm
	cfg.Auth.RememberMeTime = AuthRememberMeTime
	cfg.Auth.SignKey = []byte(AuthSignKey)
	cfg.Auth.TokenName = AuthTokenName

	cfg.DeleteShotURL.ChunkLength = 10
	cfg.DeleteShotURL.WorkersCount = 1

	return cfg
}
