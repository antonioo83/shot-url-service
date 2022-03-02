package handlers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/itchyny/base58-go"
	"math/big"
)

func GetShortURL(originalURL string, newBaseUrl string) (string, string, error) {
	urlID, err := generateShortLink(originalURL, "userId")
	return newBaseUrl + "/" + urlID, urlID, err
}

func generateShortLink(initialLink string, userID string) (string, error) {
	urlHashBytes := sha256Of(initialLink + userID)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, err := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8], err
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", errors.New("Can't encode string:" + string(bytes))
	}
	return string(encoded), nil
}
