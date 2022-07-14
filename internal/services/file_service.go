package services

import (
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/utils"
	"os"
)

func SaveToFile(fileName string, data []byte) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("i can't open a file: %w", err)
	}
	defer utils.ResourceClose(file)

	err = utils.LogErr(file.Write(data))
	if err != nil {
		return fmt.Errorf("i can't write to file: %w", err)
	}

	return nil
}
