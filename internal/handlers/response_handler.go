package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/handlers/auth/cookie"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetCreateJSONShortURLResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository) {
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
		&createShortURLs,
		func(w http.ResponseWriter, shotURLResponses []shortUrlResponse) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			jsonResponse, err := getJSONResponse("result", shotURLResponses[0].ShortUrl)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			utils.LogErr(w.Write(jsonResponse))
		},
	})
}

type shortUrlResponse struct {
	CorrelationId string `json:"correlation_id"`
	ShortUrl      string `json:"original_url"`
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

func GetCreateShortURLResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository) {
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
		&createShortURLs,
		func(w http.ResponseWriter, shotURLResponses []shortUrlResponse) {
			w.WriteHeader(http.StatusCreated)
			for _, shotURLResponse := range shotURLResponses {
				utils.LogErr(w.Write([]byte(shotURLResponse.ShortUrl)))
			}
		},
	})
}

func GetCreateShortURLBatchResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository) {
	createShortURLs, err := GetBatchRequestsFromBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getSavedShortURLResponse(savedShortURLParameters{
		w,
		r,
		config,
		repository,
		userRepository,
		createShortURLs,
		func(w http.ResponseWriter, shotURLResponses []shortUrlResponse) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			jsonResponse, err := getJSONArrayResponse(shotURLResponses)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			utils.LogErr(w.Write(jsonResponse))
		},
	})
}

func getJSONArrayResponse(shotURLResponses []shortUrlResponse) ([]byte, error) {
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
	createShortURLs *[]CreateShortURL
	responseFunc    func(w http.ResponseWriter, shotURLResponses []shortUrlResponse)
}

func getSavedShortURLResponse(p savedShortURLParameters) {
	var shortUrlResponses []shortUrlResponse
	user := getAuthUser(p.rWriter, p.request, p.userRepository)
	if isInUser, _ := p.userRepository.IsInDatabase(user.CODE); !isInUser {
		err2 := p.userRepository.Save(user)
		if err2 != nil {
			http.Error(p.rWriter, err2.Error(), http.StatusInternalServerError)
			return
		}
	}

	for _, createShortURL := range *p.createShortURLs {
		shotURL, code, err := GetShortURL(createShortURL.OriginalURL, p.request, p.config.BaseURL)
		if err != nil {
			http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		isInDB, err := p.repository.IsInDatabase(code)
		if err != nil {
			http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		if isInDB {
			shortUrlResponses = append(shortUrlResponses, shortUrlResponse{createShortURL.CorrelationId, shotURL})
			p.responseFunc(p.rWriter, shortUrlResponses)
			return
		}

		var shortURL models.ShortURL
		shortURL.Code = code
		shortURL.UserCode = user.CODE
		shortURL.CorrelationId = createShortURL.CorrelationId
		shortURL.OriginalURL = createShortURL.OriginalURL
		shortURL.ShortURL = shotURL
		err = p.repository.SaveURL(shortURL)
		if err != nil {
			http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		shortUrlResponses = append(shortUrlResponses, shortUrlResponse{shortURL.CorrelationId, shortURL.ShortURL})
	}

	p.responseFunc(p.rWriter, shortUrlResponses)
}

func getAuthUser(rWriter http.ResponseWriter, request *http.Request, userRepository interfaces.UserRepository) models.User {
	var user models.User
	token := cookie.GetToken(request)
	if token == "" || cookie.ValidateToken(token) == false {
		lastModel, _ := userRepository.GetLastModel()
		if lastModel.CODE == 0 {
			user.CODE = 1
			user.UID, _ = cookie.GenerateToken(user.CODE)
		} else {
			user = *lastModel
			user.CODE = user.CODE + 1
			user.UID, _ = cookie.GenerateToken(user.CODE)
		}
		cookie.SetToken(rWriter, user)
	} else {
		user.CODE = cookie.GetUserCode(token)
		user.UID = token
	}

	return user
}

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
	if model == nil {
		http.Error(w, "can't find model for the code: %s"+code, http.StatusNoContent)
	} else {
		w.Header().Set("Location", model.OriginalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		utils.LogErr(w.Write([]byte(model.OriginalURL)))
	}
}

func GetUserURLsResponse(w http.ResponseWriter, r *http.Request, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository) {
	user := getAuthUser(w, r, userRepository)
	models, _ := repository.FindAllByUserCode(user.CODE)
	if len(*models) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parseData := make([]map[string]interface{}, 0, 0)
	for _, model := range *models {
		var singleMap = make(map[string]interface{})
		singleMap["short_url"] = model.ShortURL
		singleMap["original_url"] = model.OriginalURL
		parseData = append(parseData, singleMap)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(parseData)
	if err != nil {
		http.Error(w, "httpStatus param is missed", http.StatusBadRequest)
	}
	utils.LogErr(w.Write(jsonResp))
}

func GetDBStatusResponse(w http.ResponseWriter, r *http.Request, databaseRepository interfaces.DatabaseRepository) {
	context := context.Background()
	conn, err := databaseRepository.Connect(context)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	databaseRepository.Disconnect(context, conn)

	w.WriteHeader(http.StatusOK)
}
