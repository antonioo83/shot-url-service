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

//GetCount gets count of short url in the storage.
func (r fileRepository) GetCount() (int, error) {
	var count int
	var model = models.ShortURL{}
	consumer, err := GetConsumer(r.filename)
	if err != nil {
		return 0, err
	}
	defer utils.ResourceClose(consumer.file)

	for consumer.scanner.Scan() {
		jsonString := consumer.scanner.Text()
		if jsonString != "" {
			err := json.Unmarshal([]byte(jsonString), &model)
			if err != nil {
				return 0, fmt.Errorf("i can't decode json request: %s", err.Error())
			}
			count++
		}
	}

	if err := consumer.scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner of a consumer got the error: %w", err)
	}

	return count, nil
}

//SaveURL saves an entity in the storage.
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

//SaveModels saves batch of entities in the storage.
func (r *fileRepository) SaveModels(models []models.ShortURL) error {
	producer, err := GetProducer(r.filename)
	if err != nil {
		return err
	}
	defer utils.ResourceClose(producer.file)

	for _, model := range models {
		err = producer.encoder.Encode(&model)
		if err != nil {
			return err
		}
	}

	return nil
}

//FindByCode finds an entity in the storage by unique code.
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

//FindAllByUserCode finds entities in the storage by unique codes.
func (r fileRepository) FindAllByUserCode(userCode int) (*map[string]models.ShortURL, error) {
	var model = models.ShortURL{}
	var models = make(map[string]models.ShortURL)
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
			if model.UserCode == userCode {
				models[model.Code] = model
			}
		}
	}

	if err := consumer.scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner of a consumer got the error: %w", err)
	}

	return &models, nil
}

//IsInDatabase check exists an entity in the storage by unique code.
func (r fileRepository) IsInDatabase(code string) (bool, error) {
	model, err := r.FindByCode(code)

	return !(model == nil), err
}

//Delete deletes entities of a user from the storage by user code.
func (r fileRepository) Delete(userCode int, codes []string) error {

	return fmt.Errorf("method wasn't implemented")
}
