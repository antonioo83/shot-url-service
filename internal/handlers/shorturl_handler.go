package handlers

import (
	"crypto/sha256"
	"fmt"
	"github.com/itchyny/base58-go"
	"math/big"
	"net/http"
	"os"
)

func GetShortUrl(originalUrl string, r *http.Request) (string, string, error) {
	urlId := generateShortLink(originalUrl, "userId")
	return "http://" + r.Host + "/" + urlId, urlId, nil
}

func generateShortLink(initialLink string, userId string) string {
	urlHashBytes := sha256Of(initialLink + userId)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8]
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}
