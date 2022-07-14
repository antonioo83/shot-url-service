package interfaces

import (
	"github.com/antonioo83/shot-url-service/internal/models"
)

type UserRepository interface {
	//Save saves a user in the storage.
	Save(model models.User) error
	//FindByCode finds a user in the storage by unique code.
	FindByCode(code int) (*models.User, error)
	//IsInDatabase check exists a user in the storage by unique code.
	IsInDatabase(code int) (bool, error)
	//GetLastModel gets a last user from the storage.
	GetLastModel() (*models.User, error)
}
