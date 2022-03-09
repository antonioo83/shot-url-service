package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/utils"
	"io"
	"net/http"
)

func GetBody(r *http.Request) (string, error) {
	defer utils.ResourceClose(r.Body)
	b, err := io.ReadAll(r.Body)

	return string(b), err
}

type shortURLRequest struct {
	URL string
}

func GetOriginalURLFromBody(r *http.Request) (string, error) {
	var request shortURLRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return "", fmt.Errorf("I can't decode json request: %w", err)
	}

	return request.URL, nil
}
