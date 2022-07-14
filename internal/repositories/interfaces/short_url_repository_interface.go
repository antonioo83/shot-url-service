package interfaces

import (
	"github.com/antonioo83/shot-url-service/internal/models"
)

type ShotURLRepository interface {
	//SaveURL saves an entity in the storage.
	SaveURL(model models.ShortURL) error
	//SaveModels saves batch of entities in the storage.
	SaveModels(models []models.ShortURL) error
	//FindByCode finds an entity in the storage by unique code.
	FindByCode(code string) (*models.ShortURL, error)
	//FindAllByUserCode finds entities in the storage by unique codes.
	FindAllByUserCode(userCode int) (*map[string]models.ShortURL, error)
	//IsInDatabase check exists an entity in the storage by unique code.
	IsInDatabase(code string) (bool, error)
	//Delete deletes entities of a user from the storage by user code.
	Delete(userCode int, codes []string) error
}
