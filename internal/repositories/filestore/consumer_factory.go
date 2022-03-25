package filestore

import (
	"bufio"
	"encoding/json"
	"os"
)

type consumer struct {
	file    *os.File
	decoder *json.Decoder
	scanner *bufio.Scanner
}

func GetConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
		scanner: bufio.NewScanner(file),
	}, nil
}
