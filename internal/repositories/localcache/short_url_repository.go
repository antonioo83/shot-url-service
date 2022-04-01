package localcache

import (
	"errors"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
)

type memoryRepository struct {
	buffer map[string]models.ShortURL
}

func NewMemoryRepository(m map[string]models.ShortURL) interfaces.ShotURLRepository {
	return &memoryRepository{m}
}

func (m *memoryRepository) SaveURL(model models.ShortURL) error {
	m.buffer[model.Code] = model

	return nil
}

func (m *memoryRepository) SaveModels(models []models.ShortURL) error {
	for _, model := range models {
		m.buffer[model.Code] = model
	}

	return nil
}

func (m *memoryRepository) FindByCode(code string) (*models.ShortURL, error) {
	model, ok := m.buffer[code]
	if !ok {
		return nil, errors.New("Can't find model in buffer for the code:" + code)
	}

	return &model, nil
}

func (m *memoryRepository) FindAllByUserCode(userCode int) (*map[string]models.ShortURL, error) {
	var models = make(map[string]models.ShortURL)
	for _, model := range m.buffer {
		if model.UserCode == userCode {
			models[model.Code] = model
		}
	}

	return &models, nil
}

func (m *memoryRepository) IsInDatabase(code string) (bool, error) {
	_, ok := m.buffer[code]

	return ok, nil
}
