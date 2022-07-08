package services

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/shot-url-service/config"
	"os"
)

// LoadConfigFile this method read a server configurations from a file in the json format.
func LoadConfigFile(configFilePath string) (*config.Config, error) {
	var configFromFile config.Config

	file, err := os.OpenFile(configFilePath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, fmt.Errorf("unable to open a configuration file: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get statistic info about configuration file: %w", err)
	}
	filesize := info.Size()
	jsonConfig := make([]byte, filesize)

	_, err = file.Read(jsonConfig)
	if err != nil {
		return nil, fmt.Errorf("i can't read a configuration file: %w", err)
	}

	err = json.Unmarshal(jsonConfig, &configFromFile)
	if err != nil {
		return nil, fmt.Errorf("i can't parse a configuration json file: %w", err)
	}

	return &configFromFile, nil
}
