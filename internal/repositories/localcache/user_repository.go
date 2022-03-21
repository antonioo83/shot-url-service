package localcache

import (
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
)

type memoryUserRepository struct {
	buffer map[int]models.User
}

func NewMemoryUserRepository(m map[int]models.User) interfaces.UserRepository {
	return &memoryUserRepository{m}
}

func (m *memoryUserRepository) Save(model models.User) error {
	m.buffer[model.CODE] = model

	return nil
}

func (m *memoryUserRepository) FindByCode(code int) (*models.User, error) {
	model, ok := m.buffer[code]
	if !ok {
		return nil, fmt.Errorf("can't find model in buffer for the code: %d", code)
	}

	return &model, nil
}

func (m *memoryUserRepository) IsInDatabase(code int) (bool, error) {
	_, ok := m.buffer[code]

	return ok, nil
}

func (m *memoryUserRepository) GetLastModel() (*models.User, error) {
	lastModel, ok := m.buffer[len(m.buffer)-1]
	if !ok {
		return &models.User{ID: 0}, nil
	}

	return &lastModel, nil
}
