package handlers

import (
	"io"
	"net/http"
)

func GetBody(r *http.Request) (string, error) {
	b, err := io.ReadAll(r.Body)

	return string(b), err
}
