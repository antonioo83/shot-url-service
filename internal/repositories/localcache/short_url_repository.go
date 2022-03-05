package localcache

import (
	"errors"
	"github.com/antonioo83/shot-url-service/internal/models"
)

var Database = make(map[string]models.ShortURL)

func SaveURL(model models.ShortURL) error {
	Database[model.Code] = model

	return nil
}

func FindByCode(code string) (*models.ShortURL, error) {
	model, ok := Database[code]
	if !ok {
		return nil, errors.New("Can't find model for the code:" + code)
	}

	return &model, nil
}

func IsHasInDatabase(code string) bool {
	_, ok := Database[code]
	if !ok {
		return false
	}

	return true
}
