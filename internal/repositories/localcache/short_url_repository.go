package localcache

import (
	"errors"
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
)

type memoryRepository struct {
	buffer map[string]models.ShortURL
}

func NewMemoryRepository(m map[string]models.ShortURL) interfaces.ShotURLRepository {
	return &memoryRepository{m}
}

//SaveURL saves an entity in the storage.
func (m *memoryRepository) SaveURL(model models.ShortURL) error {
	m.buffer[model.Code] = model

	return nil
}

//SaveModels saves batch of entities in the storage.
func (m *memoryRepository) SaveModels(models []models.ShortURL) error {
	for _, model := range models {
		m.buffer[model.Code] = model
	}

	return nil
}

//FindByCode finds an entity in the storage by unique code.
func (m *memoryRepository) FindByCode(code string) (*models.ShortURL, error) {
	model, ok := m.buffer[code]
	if !ok {
		return nil, errors.New("Can't find model in buffer for the code:" + code)
	}

	return &model, nil
}

//FindAllByUserCode finds entities in the storage by unique codes.
func (m *memoryRepository) FindAllByUserCode(userCode int) (*map[string]models.ShortURL, error) {
	var models = make(map[string]models.ShortURL)
	for _, model := range m.buffer {
		if model.UserCode == userCode {
			models[model.Code] = model
		}
	}

	return &models, nil
}

//IsInDatabase check exists an entity in the storage by unique code.
func (m *memoryRepository) IsInDatabase(code string) (bool, error) {
	_, ok := m.buffer[code]

	return ok, nil
}

//Delete deletes entities of a user from the storage by user code.
func (m *memoryRepository) Delete(userCode int, codes []string) error {

	return fmt.Errorf("method wasn't implemented")
}
