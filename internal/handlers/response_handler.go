package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antonioo83/shot-url-service/config"
	authInterfaces "github.com/antonioo83/shot-url-service/internal/handlers/auth/interfaces"
	genInterfaces "github.com/antonioo83/shot-url-service/internal/handlers/generators/interfaces"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/services"
	"github.com/antonioo83/shot-url-service/internal/utils"
	"github.com/go-chi/chi/v5"
	"net"
	"net/http"
)

type ShortURLResponse struct {
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
		func(w http.ResponseWriter, shotURLResponses []ShortURLResponse, httpStatus int) {
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
		func(w http.ResponseWriter, shotURLResponses []ShortURLResponse, httpStatus int) {
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

	getSavedShortURLResponse(savedShortURLParameters{
		w,
		r,
		config,
		repository,
		userRepository,
		userAuth,
		generator,
		createShortURLs,
		func(w http.ResponseWriter, shotURLResponses []ShortURLResponse, httpStatus int) {
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

func getJSONArrayResponse(shotURLResponses []ShortURLResponse) ([]byte, error) {
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
	responseFunc    func(w http.ResponseWriter, shotURLResponses []ShortURLResponse, httpStatus int)
}

func getSavedShortURLResponse(p savedShortURLParameters) {
	user, err := p.userAuth.GetAuthUser(p.request, p.rWriter)
	if err != nil {
		http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	var shotURLs []services.CreateShortURL
	for _, item := range *p.createShortURLs {
		shotURLs = append(shotURLs, services.CreateShortURL{OriginalURL: item.OriginalURL, CorrelationID: item.CorrelationID})
	}

	var param services.ShortURLParameters
	param.Config = p.config
	param.Repository = p.repository
	param.UserRepository = p.userRepository
	param.Generator = p.generator
	param.Host = p.request.Host
	param.User = user
	param.CreateShortURLs = &shotURLs
	result, err := services.SaveShortURLs(param)
	if err != nil {
		http.Error(p.rWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []ShortURLResponse
	for _, item := range result.ShortURLResponses {
		responses = append(responses, ShortURLResponse{CorrelationID: item.CorrelationID, ShortURL: item.ShortURL})
	}

	p.responseFunc(p.rWriter, responses, result.Status)
}

// GetOriginalURLResponse returns original URL by code.
func GetOriginalURLResponse(w http.ResponseWriter, r *http.Request, repository interfaces.ShotURLRepository) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "httpStatus param is missed", http.StatusBadRequest)
		return
	}

	result, err := services.GetShortURL(repository, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if result.Status == http.StatusTemporaryRedirect {
		w.Header().Set("Location", result.OriginalURL)
		w.WriteHeader(result.Status)
		utils.LogErr(w.Write([]byte(result.OriginalURL)))
	} else if result.Status == http.StatusGone {
		w.WriteHeader(result.Status)
		return
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

	result, err := services.GetUserShortUrls(repository, user.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result.Status == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var parseData []userURLsResponse
	for _, model := range *result.Models {
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

type statResponse struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// GetStatsResponse returns count of shot urls and count of users saved in the databases.
func GetStatsResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository,
	userRepository interfaces.UserRepository) {

	ok, err := validateIP(r, config.TrustedSubnet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	shortURLCount, err := repository.GetCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userCount, err := userRepository.GetCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var statResponse statResponse
	statResponse.URLs = shortURLCount
	statResponse.Users = userCount

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(statResponse)
	if err != nil {
		http.Error(w, "httpStatus param is missed", http.StatusBadRequest)
	}
	utils.LogErr(w.Write(jsonResp))
}

//validateIP checks IP address from the "X-Real-IP" header and CIDR is including this IP.
func validateIP(r *http.Request, subnet string) (bool, error) {
	ipStr := r.Header.Get("X-Real-IP")
	if ipStr == "" {
		return false, nil
	}

	userIP := net.ParseIP(ipStr)
	if userIP == nil {
		return false, nil
	}

	_, opNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return false, fmt.Errorf("i can't parse CIDR: %w", err)
	}

	return opNet.Contains(userIP), nil
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

// GetDeleteShortURLResponse deletes array of short URLs by array of codes.
func GetDeleteShortURLResponse(w http.ResponseWriter, r *http.Request, config config.Config, repository interfaces.ShotURLRepository,
	userAuth authInterfaces.UserAuthHandler, jobCh chan services.ShotURLDelete) {
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

	services.SendCodesForDeleteToChanel(
		jobCh,
		services.ShotURLDelete{UserCode: user.Code, Codes: *codes},
		config.DeleteShotURL.ChunkLength,
	)

	w.WriteHeader(http.StatusAccepted)
}
