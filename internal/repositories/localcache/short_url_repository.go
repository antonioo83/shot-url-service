package localcache

import (
	"errors"
	"github.com/antonioo83/shot-url-service/internal/models"
)

var database = make(map[string]models.ShortUrl)

func SaveUrl(model models.ShortUrl) bool {
	database[model.Code] = model

	return true
}

func FindByCode(code string) (*models.ShortUrl, error) {
	model, ok := database[code]
	if ok == false {
		return nil, errors.New("Can't find model for the code:" + code)
	}

	return &model, nil
}
