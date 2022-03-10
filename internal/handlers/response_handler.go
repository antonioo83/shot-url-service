package handlers

import (
	"encoding/json"
	"errors"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetCreateJSONShortURLResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository) {
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

func GetCreateShortURLResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository) {
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
		originalURL,
		func(w http.ResponseWriter, shotURL string) {
			w.WriteHeader(http.StatusCreated)
			utils.LogErr(w.Write([]byte(shotURL)))
		},
	})
}

type savedShortURLParameters struct {
	rWriter      http.ResponseWriter
	request      *http.Request
	config       config.Config
	repository   interfaces.ShotURLRepository
	originalURL  string
	responseFunc func(w http.ResponseWriter, shotURL string)
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

	var shortURL models.ShortURL
	shortURL.Code = code
	shortURL.OriginalURL = p.originalURL
	shortURL.ShortURL = shotURL
	err = p.repository.SaveURL(shortURL)
	if err != nil {
		http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	p.responseFunc(p.rWriter, shotURL)
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
