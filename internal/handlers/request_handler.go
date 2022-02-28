package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func GetJsonRequest(key string, value string) []byte {
	request := make(map[string]string)
	request[key] = value
	jsonResp, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	return jsonResp
}

func GetBody(r *http.Request) (string, error) {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)

	return string(b), err
}

type shortUrlRequest struct {
	Url string
}

func GetUrlParameter(r *http.Request) (string, error) {
	var request shortUrlRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return "", errors.New("I can't decode json request:" + err.Error())
	}

	return request.Url, nil
}
