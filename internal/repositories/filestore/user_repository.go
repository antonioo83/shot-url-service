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

//GetCount gets count of short url in the storage.
func (r userRepository) GetCount() (int, error) {
	var count int
	model := models.User{}
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

//Save saves a user in the storage.
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

//FindByCode finds a user in the storage by unique code.
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

//GetLastModel gets a last user from the storage.
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

//IsInDatabase check exists a user in the storage by unique code.
func (r userRepository) IsInDatabase(code int) (bool, error) {
	model, err := r.FindByCode(code)

	return !(model == nil), err
}
