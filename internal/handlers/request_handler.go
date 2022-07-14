package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/utils"
	"io"
	"net/http"
)

type CreateShortURL struct {
	OriginalURL   string `json:"original_url"`   // original URL
	CorrelationID string `json:"correlation_id"` // correlation ID
}

// GetBody returns CreateShortURL structure initialized by body of HTTP request.
func GetBody(r *http.Request) (*CreateShortURL, error) {
	b, err := uncompress(r)
	return &CreateShortURL{string(b), ""}, err
}

func uncompress(r *http.Request) ([]byte, error) {
	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return []byte(""), err
		}
		reader = gz
		defer utils.ResourceClose(gz)
	} else {
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return []byte(""), err
	}

	return body, nil
}

type shortURLRequest struct {
	URL string
}

// GetOriginalURLFromBody returns CreateShortURL initialized by body of HTTP request.
func GetOriginalURLFromBody(r *http.Request) (*CreateShortURL, error) {
	var request shortURLRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &CreateShortURL{request.URL, ""}, nil
}

// GetBatchRequestsFromBody returns array of CreateShortURL structures initialized by body of HTTP request.
func GetBatchRequestsFromBody(r *http.Request) (*[]CreateShortURL, error) {
	var requests []CreateShortURL
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requests)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &requests, nil
}

// GetCorrelationIDs returns array of correlation ID initialized by body of HTTP request.
func GetCorrelationIDs(r *http.Request) (*[]string, error) {
	var correlationIDs []string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&correlationIDs)
	if err != nil {
		return nil, fmt.Errorf("i can't decode json request: %w", err)
	}

	return &correlationIDs, nil
}
