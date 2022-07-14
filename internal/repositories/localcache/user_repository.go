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

//Save saves a user in the storage.
func (m *memoryUserRepository) Save(model models.User) error {
	m.buffer[model.Code] = model

	return nil
}

//FindByCode finds a user in the storage by unique code.
func (m *memoryUserRepository) FindByCode(code int) (*models.User, error) {
	model, ok := m.buffer[code]
	if !ok {
		return nil, fmt.Errorf("can't find model in buffer for the code: %d", code)
	}

	return &model, nil
}

//IsInDatabase check exists a user in the storage by unique code.
func (m *memoryUserRepository) IsInDatabase(code int) (bool, error) {
	_, ok := m.buffer[code]

	return ok, nil
}

//GetLastModel gets a last user from the storage.
func (m *memoryUserRepository) GetLastModel() (*models.User, error) {
	lastModel, ok := m.buffer[len(m.buffer)]
	if !ok {
		return &lastModel, nil
	}

	return &lastModel, nil
}
