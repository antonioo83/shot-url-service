package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func GetCreateJSONShortURLResponse(w http.ResponseWriter, shotURL string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsonResponse, err := getJSONResponse("result", shotURL)
	if err != nil {
		log.Fatal(err)
	}

	LogErr(w.Write(jsonResponse))
}

func getJSONResponse(key string, value string) ([]byte, error) {
	resp := make(map[string]string)
	resp[key] = value
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return jsonResp, errors.New("Error happened in JSON marshal")
	}

	return jsonResp, nil
}

func GetCreateShortURLResponse(w http.ResponseWriter, shotURL string) {
	w.WriteHeader(http.StatusCreated)
	LogErr(w.Write([]byte(shotURL)))
}

func GetOriginalURLResponse(w http.ResponseWriter, originalURL string) {
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	LogErr(w.Write([]byte(originalURL)))
}
