package handlers

import (
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
	originalURL, err := GetOriginalURLFromBody(r)
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
		originalURL,
		func(w http.ResponseWriter, shotURL string) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			jsonResponse, err := getJSONResponse("result", shotURL)
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

func GetCreateShortURLResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository) {
	originalURL, err := GetBody(r)
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
		originalURL,
		func(w http.ResponseWriter, shotURL string) {
			w.WriteHeader(http.StatusCreated)
			utils.LogErr(w.Write([]byte(shotURL)))
		},
	})
}

type savedShortURLParameters struct {
	rWriter        http.ResponseWriter
	request        *http.Request
	config         config.Config
	repository     interfaces.ShotURLRepository
	userRepository interfaces.UserRepository
	originalURL    string
	responseFunc   func(w http.ResponseWriter, shotURL string)
}

func getSavedShortURLResponse(p savedShortURLParameters) {
	shotURL, code, err := GetShortURL(p.originalURL, p.request, p.config.BaseURL)
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
		p.responseFunc(p.rWriter, shotURL)
		return
	}

	user := getAuthUser(p.rWriter, p.request, p.userRepository)
	if isInUser, _ := p.userRepository.IsInDatabase(user.CODE); !isInUser {
		err = p.userRepository.Save(user)
		if err != nil {
			http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var shortURL models.ShortURL
	shortURL.Code = code
	shortURL.UserCode = user.CODE
	shortURL.OriginalURL = p.originalURL
	shortURL.ShortURL = shotURL
	err = p.repository.SaveURL(shortURL)
	if err != nil {
		http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	p.responseFunc(p.rWriter, shotURL)
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
