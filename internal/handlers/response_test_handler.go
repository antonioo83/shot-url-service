package handlers

import (
	"encoding/json"
	"fmt"
)

type resultResponse struct {
	Result string
}

func GetResultParameter(body string) (string, error) {
	var response resultResponse
	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", fmt.Errorf("i can't decode json request: %w", err)
	}

	return response.Result, nil
}
