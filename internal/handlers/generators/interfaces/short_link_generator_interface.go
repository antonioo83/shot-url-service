package interfaces

import (
	"net/http"
)

type ShortLinkGenerator interface {
	GetShortURL(originalURL string, r *http.Request, newBaseURL string) (string, string, error)
}
