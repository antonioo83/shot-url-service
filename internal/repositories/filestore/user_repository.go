package filestore

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/utils"
)

type userRepository struct {
	filename string
}

func NewUserRepository(filename string) interfaces.UserRepository {
	return &userRepository{filename}
}

func (r userRepository) Save(model models.User) error {
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

func (r userRepository) FindByCode(code int) (*models.User, error) {
	model := models.User{}
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
			if model.CODE == code {
				return &model, nil
			}
		}
	}

	if err := consumer.scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner of a consumer got the error: %w", err)
	}

	return nil, nil
}

func (r userRepository) GetLastModel() (*models.User, error) {
	model := models.User{}
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
		}
	}

	if err := consumer.scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner of a consumer got the error: %w", err)
	}

	return &model, nil
}

func (r userRepository) IsInDatabase(code int) (bool, error) {
	model, err := r.FindByCode(code)

	return !(model == nil), err
}
