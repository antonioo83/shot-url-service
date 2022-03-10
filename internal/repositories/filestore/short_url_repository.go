package filestore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/utils"
)

type fileRepository struct {
	p producer
	c consumer
}

func NewFileRepository(c consumer, p producer) interfaces.ShotURLRepository {
	return &fileRepository{p, c}
}

func (r fileRepository) SaveURL(model models.ShortURL) error {
	defer utils.ResourceClose(r.p.file)
	err := r.p.encoder.Encode(&model)
	if err != nil {
		return err
	}

	return nil
}

func (r fileRepository) FindByCode(code string) (*models.ShortURL, error) {
	model := models.ShortURL{}
	defer utils.ResourceClose(r.c.file)
	for r.c.scanner.Scan() {
		jsonString := r.c.scanner.Text()
		if jsonString != "" {
			err := json.Unmarshal([]byte(jsonString), &model)
			if err != nil {
				return nil, errors.New("I can't decode json request:" + err.Error())
			}
			if model.Code == code {
				break
			}
		}
	}

	if err := r.c.scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner of a consumer got the error: %w", err)
	}

	return nil, nil
}

func (r fileRepository) IsInDatabase(code string) (bool, error) {
	model, err := r.FindByCode(code)

	return model == nil, err
}
