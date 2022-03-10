package interfaces

import (
	"github.com/antonioo83/shot-url-service/internal/models"
)

type ShotURLRepository interface {
	SaveURL(model models.ShortURL) error
	FindByCode(code string) (*models.ShortURL, error)
	IsInDatabase(code string) (bool, error)
}
