package handlers

import (
	"encoding/json"
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
