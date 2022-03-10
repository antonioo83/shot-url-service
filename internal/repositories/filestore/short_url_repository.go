package filestore

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/utils"
)

type fileRepository struct {
	filename string
}

func NewFileRepository(filename string) interfaces.ShotURLRepository {
	return &fileRepository{filename}
}

func (r fileRepository) SaveURL(model models.ShortURL) error {
	producer, err := GetProducer(r.filename)
	if err != nil {
		return err
	}
	defer utils.ResourceClose(producer.file)
	err = producer.encoder.Encode(&model)
	if err != nil {
		return err
	}

	return nil
}

func (r fileRepository) FindByCode(code string) (*models.ShortURL, error) {
	model := models.ShortURL{}
	consumer, err := GetConsumer(r.filename)
	if err != nil {
		return nil, err
	}
	defer utils.ResourceClose(consumer.file)

	for consumer.scanner.Scan() {
		jsonString := consumer.scanner.Text()
		if jsonString != "" {
			err := json.Unmarshal([]byte(jsonString), &model)
			if err != nil {
				return nil, fmt.Errorf("i can't decode json request: %s", err.Error())
			}
			if model.Code == code {
				return &model, nil
			}
		}
	}

	if err := consumer.scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner of a consumer got the error: %w", err)
	}

	return nil, nil
}

func (r fileRepository) IsInDatabase(code string) (bool, error) {
	model, err := r.FindByCode(code)

	return !(model == nil), err
}
