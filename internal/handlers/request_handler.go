package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func GetBody(r *http.Request) (string, error) {
	defer BodyClose(r.Body)
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
