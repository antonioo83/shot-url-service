package filestore

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/handlers"
	"github.com/antonioo83/shot-url-service/internal/models"
	"log"
	"os"
)

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func SaveURL(model models.ShortURL, config config.Config) error {
	if config.IsUseFileStore == false {

		return nil
	}

	p, err := getProducer(config.FileStoragePath)
	if err != nil {
		return err
	}
	defer handlers.FileClose(p.file)

	err = p.encoder.Encode(&model)
	if err != nil {
		return err
	}

	return nil
}

func getProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
	scanner *bufio.Scanner
}

func LoadModels(database map[string]models.ShortURL, model models.ShortURL, config config.Config) (map[string]models.ShortURL, error) {
	if config.IsUseFileStore == false {

		return database, nil
	}

	consumer, err := getConsumer(config.FileStoragePath)
	if err != nil {
		return nil, err
	}
	defer handlers.FileClose(consumer.file)

	for consumer.scanner.Scan() {
		jsonString := consumer.scanner.Text()
		if jsonString != "" {
			err := json.Unmarshal([]byte(jsonString), &model)
			if err != nil {
				return nil, errors.New("I can't decode json request:" + err.Error())
			}
			database[model.Code] = model
		}
	}

	if err := consumer.scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return database, nil
}

func getConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *consumer) ReadEvent() (*models.ShortURL, error) {
	// одиночное сканирование до следующей строки
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	// читаем данные из scanner
	data := c.scanner.Bytes()

	model := models.ShortURL{}
	err := json.Unmarshal(data, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
