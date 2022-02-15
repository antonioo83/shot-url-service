package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type myData struct {
	Url string
}

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
	var data myData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		return "", errors.New("I can't decode json request:" + err.Error())
	}

	return data.Url, nil
}
