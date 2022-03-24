package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetBody(r *http.Request) (*CreateShortURL, error) {
	b, err := uncompress(r)
	return &CreateShortURL{string(b), ""}, err
}

type shortURLRequest struct {
	URL string
}

type CreateShortURL struct {
	OriginalURL   string `json:"original_url"`
	CorrelationId string `json:"correlation_id"`
}

func GetOriginalURLFromBody(r *http.Request) (*CreateShortURL, error) {
	var request shortURLRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &CreateShortURL{request.URL, ""}, nil
}

func GetBatchRequestsFromBody(r *http.Request) (*[]CreateShortURL, error) {
	var requests []CreateShortURL
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requests)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &requests, nil
}

func uncompress(r *http.Request) ([]byte, error) {
	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return []byte("test"), err
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return []byte(""), err
	}

	return body, nil
}
