package localcache

import (
	"errors"
	"github.com/antonioo83/shot-url-service/internal/models"
)

var database = make(map[string]models.ShortURL)

func SaveURL(model models.ShortURL) bool {
	database[model.Code] = model

	return true
}

func FindByCode(code string) (*models.ShortURL, error) {
	model, ok := database[code]
	if !ok {
		return nil, errors.New("Can't find model for the code:" + code)
	}

	return &model, nil
}
