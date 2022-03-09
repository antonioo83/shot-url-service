package utils

import (
	"io"
	"log"
)

func LogErr(n int, err error) int {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}

	return n
}

func ResourceClose(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Printf("Can't close resource: %v", err)
	}
}
