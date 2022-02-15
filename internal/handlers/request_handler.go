package handlers

import (
	"errors"
	"io"
	"net/http"
)

func GetQuery(name string, r *http.Request) (string, error) {
	parameter := ""
	if r.Method == http.MethodGet {
		parameter = r.URL.Query().Get(name)
		if parameter == "" {
			return "", errors.New("The query '" + name + "' parameter is missing")
		}
	}

	return parameter, nil
}

func GetUrlParameter(r *http.Request) (string, error) {
	b, err := io.ReadAll(r.Body)

	return string(b), err
}
