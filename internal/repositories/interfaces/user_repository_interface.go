package interfaces

import (
	"github.com/antonioo83/shot-url-service/internal/models"
)

type UserRepository interface {
	Save(model models.User) error
	FindByCode(code int) (*models.User, error)
	IsInDatabase(code int) (bool, error)
	GetLastModel() (*models.User, error)
}
