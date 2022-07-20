package services

import (
	"errors"
	"github.com/antonioo83/shot-url-service/config"
	genInterfaces "github.com/antonioo83/shot-url-service/internal/handlers/generators/interfaces"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"net/http"
)

type ShortURLParameters struct {
	Config          config.Config
	Repository      interfaces.ShotURLRepository
	UserRepository  interfaces.UserRepository
	Generator       genInterfaces.ShortLinkGenerator
	Host            string
	User            models.User
	CreateShortURLs *[]CreateShortURL
}

type CreateShortURL struct {
	OriginalURL   string `json:"original_url"`   // original URL
	CorrelationID string `json:"correlation_id"` // correlation ID
}

type ShortURLResult struct {
	Status            int
	ShortURLResponses []shortURLResponse
}

type shortURLResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func SaveShortURLs(p ShortURLParameters) (*ShortURLResult, error) {
	var shortURLResponses []shortURLResponse
	var shortURLModels []models.ShortURL
	for _, createShortURL := range *p.CreateShortURLs {
		shotURL, code, err := p.Generator.GetShortURL(createShortURL.OriginalURL, p.Host, p.Config.BaseURL)
		if err != nil {
			return nil, err
		}

		isInUser, err := p.UserRepository.IsInDatabase(p.User.Code)
		if err != nil {
			return nil, err
		}
		if !isInUser {
			err = p.UserRepository.Save(p.User)
			if err != nil {
				return nil, err
			}
		}

		isInDB, err := p.Repository.IsInDatabase(code)
		if err != nil {
			return nil, err
		}

		var shortURL models.ShortURL
		shortURL.Code = code
		shortURL.UserCode = p.User.Code
		shortURL.CorrelationID = createShortURL.CorrelationID
		shortURL.OriginalURL = createShortURL.OriginalURL
		shortURL.ShortURL = shotURL
		shortURL.Active = true
		if p.Config.IsUseDatabase || !isInDB {
			shortURLModels = append(shortURLModels, shortURL)
		}
		shortURLResponses = append(shortURLResponses, shortURLResponse{shortURL.CorrelationID, shortURL.ShortURL})
	}

	err := p.Repository.SaveModels(shortURLModels)
	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == pgerrcode.UniqueViolation {
			return &ShortURLResult{Status: http.StatusConflict, ShortURLResponses: shortURLResponses}, nil
		} else {
			return nil, err
		}
	}

	return &ShortURLResult{Status: http.StatusCreated, ShortURLResponses: shortURLResponses}, nil
}

type ShortURLGetResult struct {
	Status      int
	OriginalURL string
}

func GetShortURL(repository interfaces.ShotURLRepository, code string) (*ShortURLGetResult, error) {
	model, err := repository.FindByCode(code)
	if err != nil {
		return nil, err
	}
	if !model.Active {
		return &ShortURLGetResult{Status: http.StatusGone, OriginalURL: ""}, nil
	}
	if model == nil {
		return &ShortURLGetResult{Status: http.StatusNoContent, OriginalURL: ""}, nil
	}

	return &ShortURLGetResult{Status: http.StatusTemporaryRedirect, OriginalURL: model.OriginalURL}, nil
}

type UserShortURLGetResult struct {
	Status int
	Models *map[string]models.ShortURL
}

func GetUserShortUrls(repository interfaces.ShotURLRepository, userCode int) (*UserShortURLGetResult, error) {
	var response UserShortURLGetResult
	models, err := repository.FindAllByUserCode(userCode)
	if err != nil {
		return &response, err
	}
	if len(*models) == 0 {
		return &UserShortURLGetResult{Status: http.StatusNoContent}, nil
	}

	return &UserShortURLGetResult{Status: http.StatusOK, Models: models}, nil
}

type ShotURLDelete struct {
	UserCode int
	Codes    []string
}

func SendCodesForDeleteToChanel(jobCh chan ShotURLDelete, shortURLDelete ShotURLDelete, chunkLength int) {
	var chunkCodes []string
	for _, code := range shortURLDelete.Codes {
		chunkCodes = append(chunkCodes, code)
		if len(chunkCodes) == chunkLength {
			jobCh <- ShotURLDelete{UserCode: shortURLDelete.UserCode, Codes: chunkCodes}
			chunkCodes = []string{}
		}
	}
	if len(chunkCodes) > 0 {
		jobCh <- ShotURLDelete{UserCode: shortURLDelete.UserCode, Codes: chunkCodes}
	}
}

// RunDeleteShortURLWorker Run workers to delete shot URLs from database.
func RunDeleteShortURLWorker(jobCh chan ShotURLDelete, repository interfaces.ShotURLRepository, workersCount int) {
	for i := 0; i < workersCount; i++ {
		go func() {
			for shotURLDelete := range jobCh {
				repository.Delete(shotURLDelete.UserCode, shotURLDelete.Codes)
			}
		}()
	}
}
