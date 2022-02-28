package handlers

import (
	"encoding/json"
	"errors"
	"log"
)

func GetJsonResponse(key string, value string) []byte {
	resp := make(map[string]string)
	resp[key] = value
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	return jsonResp
}

type resultResponse struct {
	Result string
}

func GetResultParameter(body string) (string, error) {
	var response resultResponse
	err := json.Unmarshal([]byte(body), &response)
	//err := decoder.Decode(&response)
	if err != nil {
		return "", errors.New("I can't decode json request:" + err.Error())
	}

	return response.Result, nil
}
