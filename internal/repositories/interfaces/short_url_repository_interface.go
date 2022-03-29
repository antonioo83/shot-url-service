package interfaces

import (
	"github.com/antonioo83/shot-url-service/internal/models"
)

type ShotURLRepository interface {
	SaveURL(model models.ShortURL) error
	SaveModels(models map[int]models.ShortURL) error
	FindByCode(code string) (*models.ShortURL, error)
	FindAllByUserCode(userCode int) (*map[string]models.ShortURL, error)
	IsInDatabase(code string) (bool, error)
}
