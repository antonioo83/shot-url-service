package generators

import (
	"crypto/sha256"
	"errors"
	"fmt"
	interfaces2 "github.com/antonioo83/shot-url-service/internal/handlers/generators/interfaces"
	"github.com/itchyny/base58-go"
	"math/big"
	"net/http"
)

type shortLinkGenerator struct {
}

// NewShortLinkDefaultGenerator create userAuthHandler instance.
func NewShortLinkDefaultGenerator() interfaces2.ShortLinkGenerator {
	return &shortLinkGenerator{}
}

// GetShortURL returns generated short URL.
func (a shortLinkGenerator) GetShortURL(originalURL string, r *http.Request, newBaseURL string) (string, string, error) {
	urlID, err := a.generateShortLink(originalURL, "userId")

	if newBaseURL == "" {
		return "http://" + r.Host + "/" + urlID, urlID, err
	}

	return newBaseURL + "/" + urlID, urlID, err
}

//Review ! rktkov: Может спрятать реализацию кодирования тоже за интерфейс и в отдельный пакет вынести?
func (a shortLinkGenerator) generateShortLink(initialLink string, userID string) (string, error) {
	urlHashBytes := a.sha256Of(initialLink + userID)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, err := a.base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8], err
}

func (a shortLinkGenerator) sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func (a shortLinkGenerator) base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", errors.New("Can't encode string:" + string(bytes))
	}
	return string(encoded), nil
}
