package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetBody(r *http.Request) (string, error) {
	b, err := uncompress(r)
	//defer utils.ResourceClose(r.Body)
	//b, err = io.ReadAll(r.Body)

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
		return "", fmt.Errorf("i can't decode json request: %w", err)
	}

	return request.URL, nil
}

func uncompress(r *http.Request) ([]byte, error) {
	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return []byte(""), err
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
