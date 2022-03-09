package handlers

import (
	"encoding/json"
	"fmt"
)

func GetJSONRequest(key string, value string) ([]byte, error) {
	request := make(map[string]string)
	request[key] = value
	jsonResp, err := json.Marshal(request)
	if err != nil {
		return []byte(""), fmt.Errorf("I can't decode json request: %w", err)
	}

	return jsonResp, nil
}
