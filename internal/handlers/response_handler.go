package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/antonioo83/shot-url-service/config"
	authInterfaces "github.com/antonioo83/shot-url-service/internal/handlers/auth/interfaces"
	genInterfaces "github.com/antonioo83/shot-url-service/internal/handlers/generators/interfaces"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"net/http"
)

type shortURLResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// GetCreateJSONShortURLResponse creates a short URL by json request in the storage and returns the response.
func GetCreateJSONShortURLResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository,
	userRepository interfaces.UserRepository, userAuth authInterfaces.UserAuthHandler, generator genInterfaces.ShortLinkGenerator) {
	createShortURL, err := GetOriginalURLFromBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var createShortURLs []CreateShortURL
	createShortURLs = append(createShortURLs, *createShortURL)
	getSavedShortURLResponse(savedShortURLParameters{
		w,
		r,
		config,
		repository,
		userRepository,
		userAuth,
		generator,
		&createShortURLs,
		func(w http.ResponseWriter, shotURLResponses []shortURLResponse, httpStatus int) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(httpStatus)
			jsonResponse, err := getJSONResponse("result", shotURLResponses[0].ShortURL)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			utils.LogErr(w.Write(jsonResponse))
		},
	})
}

func getJSONResponse(key string, value string) ([]byte, error) {
	resp := make(map[string]string)
	resp[key] = value
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return jsonResp, errors.New("error happened in JSON marshal")
	}

	return jsonResp, nil
}

// GetCreateShortURLResponse creates a short URL in the storage and returns the response.
func GetCreateShortURLResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository,
	userRepository interfaces.UserRepository, userAuth authInterfaces.UserAuthHandler, generator genInterfaces.ShortLinkGenerator) {
	createShortURL, err := GetBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var createShortURLs []CreateShortURL
	createShortURLs = append(createShortURLs, *createShortURL)
	getSavedShortURLResponse(savedShortURLParameters{
		w,
		r,
		config,
		repository,
		userRepository,
		userAuth,
		generator,
		&createShortURLs,
		func(w http.ResponseWriter, shotURLResponses []shortURLResponse, httpStatus int) {
			w.WriteHeader(httpStatus)
			utils.LogErr(w.Write([]byte(shotURLResponses[0].ShortURL)))
		},
	})
}

// GetCreateShortURLBatchResponse creates a array of short URLs in the storage and returns the response.
func GetCreateShortURLBatchResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository,
	userRepository interfaces.UserRepository, userAuth authInterfaces.UserAuthHandler, generator genInterfaces.ShortLinkGenerator) {
	createShortURLs, err := GetBatchRequestsFromBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//I couldn't move function to service/app package because golang was getting an error with message: import cycle not allowed. I need consultation.
	getSavedShortURLResponse(savedShortURLParameters{
		w,
		r,
		config,
		repository,
		userRepository,
		userAuth,
		generator,
		createShortURLs,
		func(w http.ResponseWriter, shotURLResponses []shortURLResponse, httpStatus int) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(httpStatus)
			jsonResponse, err := getJSONArrayResponse(shotURLResponses)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			utils.LogErr(w.Write(jsonResponse))
		},
	})
}

func getJSONArrayResponse(shotURLResponses []shortURLResponse) ([]byte, error) {
	jsonResp, err := json.Marshal(shotURLResponses)
	if err != nil {
		return jsonResp, errors.New("error happened in JSON marshal")
	}

	return jsonResp, nil
}

type savedShortURLParameters struct {
	rWriter         http.ResponseWriter
	request         *http.Request
	config          config.Config
	repository      interfaces.ShotURLRepository
	userRepository  interfaces.UserRepository
	userAuth        authInterfaces.UserAuthHandler
	generator       genInterfaces.ShortLinkGenerator
	createShortURLs *[]CreateShortURL
	responseFunc    func(w http.ResponseWriter, shotURLResponses []shortURLResponse, httpStatus int)
}

func getSavedShortURLResponse(p savedShortURLParameters) {
	var shortURLResponses []shortURLResponse
	user, err := p.userAuth.GetAuthUser(p.request, p.rWriter)
	if err != nil {
		http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	var shortURLModels []models.ShortURL
	for _, createShortURL := range *p.createShortURLs {
		shotURL, code, err := p.generator.GetShortURL(createShortURL.OriginalURL, p.request, p.config.BaseURL)
		if err != nil {
			http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		isInUser, err := p.userRepository.IsInDatabase(user.Code)
		if err != nil {
			http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		if !isInUser {
			err = p.userRepository.Save(user)
			if err != nil {
				http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		isInDB, err := p.repository.IsInDatabase(code)
		if err != nil {
			http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		var shortURL models.ShortURL
		shortURL.Code = code
		shortURL.UserCode = user.Code
		shortURL.CorrelationID = createShortURL.CorrelationID
		shortURL.OriginalURL = createShortURL.OriginalURL
		shortURL.ShortURL = shotURL
		shortURL.Active = true
		if p.config.IsUseDatabase || !isInDB {
			shortURLModels = append(shortURLModels, shortURL)
		}
		shortURLResponses = append(shortURLResponses, shortURLResponse{shortURL.CorrelationID, shortURL.ShortURL})
	}

	err = p.repository.SaveModels(shortURLModels)
	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == pgerrcode.UniqueViolation {
			p.responseFunc(p.rWriter, shortURLResponses, http.StatusConflict)
			return
		} else {
			http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	p.responseFunc(p.rWriter, shortURLResponses, http.StatusCreated)
}

// GetOriginalURLResponse returns original URL by code.
func GetOriginalURLResponse(w http.ResponseWriter, r *http.Request, repository interfaces.ShotURLRepository) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "httpStatus param is missed", http.StatusBadRequest)
		return
	}
	model, err := repository.FindByCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !model.Active {
		w.WriteHeader(http.StatusGone)
		return
	}
	// I can't remove else expression because it'll brake a wet test.
	if model == nil {
		http.Error(w, "can't find model for the code: %s"+code, http.StatusNoContent)
	} else {
		w.Header().Set("Location", model.OriginalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		utils.LogErr(w.Write([]byte(model.OriginalURL)))
	}
}

type userURLsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// GetUserURLsResponse returns all short URLs for a given user.
func GetUserURLsResponse(w http.ResponseWriter, r *http.Request, repository interfaces.ShotURLRepository,
	userRepository interfaces.UserRepository, userAuth authInterfaces.UserAuthHandler) {
	user, err := userAuth.GetAuthUser(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	models, err := repository.FindAllByUserCode(user.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(*models) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var parseData []userURLsResponse
	for _, model := range *models {
		parseData = append(parseData, userURLsResponse{model.ShortURL, model.OriginalURL})
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(parseData)
	if err != nil {
		http.Error(w, "httpStatus param is missed", http.StatusBadRequest)
	}
	utils.LogErr(w.Write(jsonResp))
}

// GetDBStatusResponse returns database status.
func GetDBStatusResponse(w http.ResponseWriter, databaseRepository interfaces.DatabaseRepository) {
	context := context.Background()
	conn, err := databaseRepository.Connect(context)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = conn.Ping(context)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	databaseRepository.Disconnect(context, conn)

	w.WriteHeader(http.StatusOK)
}

type ShotURLDelete struct {
	UserCode int
	Codes    []string
}

// GetDeleteShortURLResponse deletes array of short URLs by array of codes.
func GetDeleteShortURLResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository,
	userAuth authInterfaces.UserAuthHandler, jobCh chan ShotURLDelete) {
	user, err := userAuth.GetAuthUser(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	codes, err := GetCorrelationIDs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendCodesForDeleteToChanel(
		jobCh,
		ShotURLDelete{user.Code, *codes},
		config.DeleteShotURL.ChunkLength,
	)

	w.WriteHeader(http.StatusAccepted)
}

// I am using the workerpool pattern because I don't see reasons in using the fanIn pattern.
func sendCodesForDeleteToChanel(jobCh chan ShotURLDelete, shortURLDelete ShotURLDelete, chunkLength int) {
	var chunkCodes []string
	for _, code := range shortURLDelete.Codes {
		chunkCodes = append(chunkCodes, code)
		if len(chunkCodes) == chunkLength {
			jobCh <- ShotURLDelete{shortURLDelete.UserCode, chunkCodes}
			chunkCodes = []string{}
		}
	}
	if len(chunkCodes) > 0 {
		jobCh <- ShotURLDelete{shortURLDelete.UserCode, chunkCodes}
	}
}
