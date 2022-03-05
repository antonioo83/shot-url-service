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

type shortURLRequest struct {
	Url string
}

func GetURLParameter(r *http.Request) (string, error) {
	var request shortURLRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return "", errors.New("I can't decode json request:" + err.Error())
	}

	return request.Url, nil
}
