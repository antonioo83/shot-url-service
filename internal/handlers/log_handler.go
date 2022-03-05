package handlers

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

func BodyClose(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Fatal(err)
	}
}

func FileClose(file io.Closer) {
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}
