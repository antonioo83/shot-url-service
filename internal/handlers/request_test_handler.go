package handlers

import (
	"encoding/json"
	"log"
)

func GetJSONRequest(key string, value string) []byte {
	request := make(map[string]string)
	request[key] = value
	jsonResp, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	return jsonResp
}
