package handlers

import (
	"encoding/json"
	"errors"
)

type resultResponse struct {
	Result string
}

func GetResultParameter(body string) (string, error) {
	var response resultResponse
	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", errors.New("I can't decode json request:" + err.Error())
	}

	return response.Result, nil
}
